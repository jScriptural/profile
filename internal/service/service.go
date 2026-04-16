package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"profile/internal/models"
	"log"
	"sync"
	"time"
	"fmt"
	"errors"
)

type Store interface {
	SaveProfile(ctx context.Context, profile *models.Profile) error

	GetProfileByName(ctx context.Context, name string) (*models.Profile, error)

	GetProfileByID(ctx context.Context, id string) (*models.Profile, error)

	GetProfile(ctx context.Context, gender, countryID, ageGroup string) ([]*models.Profile, error)

	DeleteProfileByID(ctx context.Context, id string) error

}

type Service struct {
	store  Store
	client *http.Client
}

func NewService(s Store) *Service {
	return &Service{
		store:  s,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *Service) GetOrCreateProfile(ctx context.Context, name string) (*models.Profile, bool, error) {
	p, err := s.store.GetProfileByName(ctx, name)
	if err == nil {
		return p, false, nil
	}

	if err != nil && !errors.Is(err, models.ErrNoRows) {
		log.Printf("GetorCreateProfile: %v",err)
		return nil, false, err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var (
		gRes models.GenderizeResponse
		aRes models.AgifyResponse
		nRes models.NationalizeResponse
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := s.fetchGender(ctx, name, &gRes); err != nil {
			errChan <- fmt.Errorf("genderize: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.fetchAge(ctx, name, &aRes); err != nil {
			errChan <- fmt.Errorf("agify: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := s.fetchNation(ctx, name, &nRes); err != nil {
			errChan <- fmt.Errorf("nationalize: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			log.Printf("GetorCreateProfile: %v",err)
			return nil, false, err
		}
	}

	if len(nRes.Country) == 0 || aRes.Age == nil || gRes.Gender == nil || gRes.Count == 0 {
		return nil, false, fmt.Errorf("GetOrCreateProfile: %w",models.Err502);
	}

	prof := s.assembleProfile(name, &gRes, &aRes, &nRes)

	if err := s.store.SaveProfile(ctx, prof); err != nil {
		log.Printf("GetOrCreateProfile: %v",err)
		return nil, false, fmt.Errorf("GetOrCreateProfile: %w",err);
	}

	return prof, true, nil
}



func (s *Service) RetrieveProfileByID(ctx context.Context, id string) (*models.Profile, error) {

	p,err := s.store.GetProfileByID(ctx,id);
	if err != nil {
		return nil,fmt.Errorf("RetrieveProfileByID: %w",err);
	}

	return p,nil;
}











/****************************************
*                                       *
*            HELPER FUNCS               *
*                                       *
*****************************************/



func (s *Service) fetchNation(ctx context.Context, name string, nRes *models.NationalizeResponse) error {
	u, err := url.Parse("https://api.nationalize.io")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("name", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nationalize api returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(nRes); err != nil {
		return fmt.Errorf("nationalize: %w", err)
	}

	return nil
}

func (s *Service) fetchGender(ctx context.Context, name string, gRes *models.GenderizeResponse) error {
	u, err := url.Parse("https://api.genderize.io")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("name", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("genderize api returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(gRes); err != nil {
		return fmt.Errorf("genderize: %w", err)
	}

	return nil
}

func (s *Service) fetchAge(ctx context.Context, name string, aRes *models.AgifyResponse) error {

	u, err := url.Parse("https://api.agify.io")
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("name", name)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("agify api returned status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(aRes); err != nil {
		return fmt.Errorf("agify: %w", err)
	}

	return nil
}

func (s *Service) assembleProfile(name string, gRes *models.GenderizeResponse, aRes *models.AgifyResponse, nRes *models.NationalizeResponse) *models.Profile {
	p := models.Profile{}

	p.ID,_= uuid.NewV7()
	p.Name = name
	p.Gender = *gRes.Gender
	p.GenderProbability = gRes.Probability
	p.SampleSize = int(gRes.Count)
	p.Age = *aRes.Age
	p.AgeGroup = aRes.GetAgeGroup()

	c := nRes.GetMostProbableNationality()
	p.CountryID = c.CountryID
	p.CountryProbability = c.Probability
	p.CreatedAt = time.Now().UTC()

	return &p
}
