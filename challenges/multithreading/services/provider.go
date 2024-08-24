package services

import "challenges_multithreading/models"

type CepProvider interface {
	GetName() string
	FindCep(string) (*models.CepResult, error)
}
