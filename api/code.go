package api

/*
   @File: code.go
   @Author: khaosles
   @Time: 2023/8/17 21:02
   @Desc:
*/

const (
	StatusOK = 20000 // RFC 9110, 15.3.1

	StatusBadRequest       = 40000 // RFC 9110, 15.5.1
	StatusParamError       = 40001 // RFC 9110, 15.5.1
	StatusUnauthorized     = 40100 // RFC 9110, 15.5.2
	StatusNotToken         = 40101 // RFC 9110, 15.5.2
	StatusTokenExpired     = 40102 // RFC 9110, 15.5.2
	StatusTokenIpError     = 40103 // RFC 9110, 15.5.2
	StatusForbidden        = 40300 // RFC 9110, 15.5.4
	StatusNotFound         = 40400 // RFC 9110, 15.5.5
	StatusMethodNotAllowed = 40500 // RFC 9110, 15.5.6

	StatusInternalServerError = 50000 // RFC 9110, 15.6.1
	StatusNotImplemented      = 50100 // RFC 9110, 15.6.2
	StatusBadGateway          = 50200 // RFC 9110, 15.6.3
)
