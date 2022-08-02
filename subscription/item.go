package main

type Card struct {
	Id          uint   `json:"id"`
	Uid         string `json:"uid"`
	ValidCard   string `json:"valid_card"`
	Token       string `json:"token"`
	InvalidCard string `json:"invalid_card"`
	Month       string `json:"month"`
	Year        string `json:"year"`
	CCV         string `json:"ccv"`
	CCVAmex     string `json:"ccv_amex"`
}
