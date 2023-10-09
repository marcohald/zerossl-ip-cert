/*
 * Copyright [2022] [tinkernels (github.com/tinkernels)]
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zerosslIPCert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	//	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Client is a client for ZeroSSL.
// Refer: https://zerossl.com/documentation/api
type Client struct {
	ApiKey string // API key
}

// GetCert returns a certificate.
func (c *Client) GetCert(id string) (cert CertificateInfoModel, err error) {
	req_ := ApiReqFactory.GetCertificate(c.ApiKey, id)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		return CertificateInfoModel{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		return CertificateInfoModel{}, fmt.Errorf("ZeroSSL API returned status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&cert)
	if err != nil {
		if &cert == nil {
			return
		}
		log.Println(err)
		// The validation field in api response can an empty array, using the partially unmarshalled value.
		return cert, nil
	}
	return
}

// CreateCert creates a certificate with the given parameters.
func (c *Client) CreateCert(domains, csr, days, isStrictDomains string) (cert CertificateInfoModel, err error) {
	req_ := ApiReqFactory.CreateCertificate(c.ApiKey, domains, csr, days, isStrictDomains)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		log.Println(err)
		return CertificateInfoModel{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		return CertificateInfoModel{}, fmt.Errorf("ZeroSSL API returned status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&cert)
	if err != nil {
		return CertificateInfoModel{}, err
	}
	return
}

// DeleteCert deletes a certificate.
func (c *Client) DeleteCert(id string) (err error) {
	req_ := ApiReqFactory.CancelCertificate(c.ApiKey, id)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %s", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("ZeroSSL API returned status code %d, body %s", resp.StatusCode, string(body))
	}
	return
}

// VerifyDomains verifies domains of specified certificate with given validation info.
func (c *Client) VerifyDomains(certID, validationMethod, validationEmail string) (verifyDomainsRsp VerifyDomainsModel, err error) {
	req_ := ApiReqFactory.VerifyDomains(c.ApiKey, certID, validationMethod, validationEmail)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		return VerifyDomainsModel{}, fmt.Errorf("ZeroSSL API returned status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&verifyDomainsRsp)
	if err != nil {
		return VerifyDomainsModel{}, err
	}
	return
}

// VerificationStatus returns the verification status of a certificate.
func (c *Client) VerificationStatus(certID string) (verificationStatusRsp VerificationStatusModel, err error) {
	req_ := ApiReqFactory.VerificationStatus(c.ApiKey, certID)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		return VerificationStatusModel{}, fmt.Errorf("ZeroSSL API returned status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&verificationStatusRsp)
	if err != nil {
		return VerificationStatusModel{}, err
	}
	return
}

// DownloadCertInline returns the certificate in PEM format.
func (c *Client) DownloadCertInline(certID, includeCrossSigned string) (cert CertificateContentModel, err error) {
	req_ := ApiReqFactory.DownloadCertificateInline(c.ApiKey, certID, includeCrossSigned)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		return CertificateContentModel{}, fmt.Errorf("ZeroSSL API returned status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(&cert)
	if err != nil {
		return CertificateContentModel{}, err
	}
	return
}

// ListCerts returns a list of certificates with optional filters.
func (c *Client) ListCerts(status, search, limit, page string) (listCertsRsp ListCertsModel, err error) {
	log.Printf("status: %s search: %s limit: %s page: %s", status, search, limit, page)
	req_ := ApiReqFactory.ListCertificates(c.ApiKey, status, search, limit, page)
	resp, err := http.DefaultClient.Do(req_)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode >= 400 {
		return ListCertsModel{}, fmt.Errorf("ZeroSSL API returned status code %d, boday %s", resp.StatusCode, resp.Body)
	}
	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)
	//body, err := io.ReadAll(&buf)
	body, err := io.ReadAll(tee)
	log.Println(string(body))
	err = json.NewDecoder(&buf).Decode(&listCertsRsp)

	if err != nil {
		if &listCertsRsp == nil {
			return
		}
		// Read the response body
		//body, err := io.ReadAll(test)
		//log.Println(string(body))
		log.Println(err)
		// The validation field in api response can an empty array, using the partially unmarshalled value.
		return listCertsRsp, nil
	}
	return
}

func (c *Client) CleanUnfinished() (err error) {
	log.Println("Cleaning unfinished certificates")
	perPage_ := 100
	page_ := 1
	//certs, err := c.ListCerts("", "draft,pending_validation", strconv.Itoa(perPage_), strconv.Itoa(page_))
	certs, err := c.ListCerts("", "", strconv.Itoa(perPage_), strconv.Itoa(page_))
	for i := 0; page_-1 <= certs.TotalCount/perPage_; page_++ {
		certs, err := c.ListCerts("", "", strconv.Itoa(perPage_), strconv.Itoa(page_))
		log.Printf("i: %d certs.TotalCount/perPage_: %d", i, certs.TotalCount/perPage_)
		log.Printf("page: %d perPage: %d ResultCount: %d TotalCount: %d", page_, perPage_, certs.ResultCount, certs.TotalCount)
		if err != nil {
			log.Println(err)
			break
		}

		for _, cert := range certs.Results {
			// Cleaning up certificates that are not finished (including draft,PendingValidation).
			//log.Printf("Loop Name %s in %s status, id %s", cert.CommonName, cert.Status, cert.ID)
			if cert.Status == CertStatus.Draft || cert.Status == CertStatus.PendingValidation {
				log.Printf("Cleaning %s in %s status, id %s", cert.CommonName, cert.Status, cert.ID)
				err = c.DeleteCert(cert.ID)
				if err != nil {
					log.Println(err)
				}
			}
		}
		log.Printf("certs.ResultCount: %d certs.Limit: %d", certs.ResultCount, certs.Limit)
	}
	return
}
