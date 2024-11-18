package ociclient

// import (
// 	"context"
// 	"encoding/base64"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	"oras.land/oras-go/v2/registry/remote/credentials"
// )

// type Credentials struct {
// 	Username string
// 	Password string
// 	encoded  string
// }

// type Token struct {
// 	AccessToken string `json:"access_token"`
// }

// func parseRealmHeader(header string) (realm, service string, err error) {
// 	parts := strings.Split(header, ",")
// 	for _, part := range parts {
// 		part = strings.TrimSpace(part)
// 		if strings.HasPrefix(part, "realm=") {
// 			realm = strings.Trim(part[len("realm="):], "\"")
// 		} else if strings.HasPrefix(part, "service=") {
// 			service = strings.Trim(part[len("service="):], "\"")
// 		}
// 	}

// 	if realm == "" {
// 		return "", "", errors.New("missing realm in WWW-Authenticate header")
// 	}

// 	return realm, service, nil
// }

// // All OCI Registry API requests should support basic auth as part of the spec
// func (c *Client) SetBasicAuth(username, password string) {
// 	userpass := fmt.Sprintf("%s:%s", username, password)
// 	encoded := base64.StdEncoding.EncodeToString([]byte(userpass))
// 	authHeader := fmt.Sprintf("Basic %s", encoded)
// 	c.Credentials = &Credentials{
// 		Username: username,
// 		Password: password,
// 		encoded:  authHeader,
// 	}
// }

// func (c *Client) SetBearerAuth(token string) {
// 	authHeader := fmt.Sprintf("Bearer %s", token)
// 	c.Credentials.encoded = authHeader
// }

// func (c *Client) NewAuth(username string, password string, host string) error {
// 	c.SetBasicAuth(username, password)
// 	endpoint := fmt.Sprintf("https://%s/v2/", host)

// 	req, err := http.NewRequest("GET", endpoint, nil)
// 	if err != nil {
// 		return err
// 	}
// 	if c.Credentials != nil {
// 		req.Header.Set("Authorization", c.Credentials.encoded)
// 	}

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusUnauthorized {
// 		authHeader := resp.Header.Get("WWW-Authenticate")
// 		if authHeader == "" {
// 			return errors.New("no WWW-Authenticate header found")
// 		}

// 		realm, service, err := parseRealmHeader(authHeader)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println("-- DEBUG: Checking response status code for realm: ", resp.StatusCode, realm, service)
// 		if err = c.SetRealm(host, realm, service); err != nil {
// 			return err
// 		}
// 	} else {
// 		fmt.Printf("-- DEBUG:%d", resp.StatusCode)
// 	}

// 	return nil
// }

// func (c *Client) SetRealm(host string, realm string, service string) error {
// 	req, err := http.NewRequest("GET", realm, nil)
// 	if err != nil {
// 		return err
// 	}

// 	if c.Credentials != nil {
// 		req.Header.Set("Authorization", c.Credentials.encoded)
// 	}

// 	q := req.URL.Query()
// 	q.Add("service", service)
// 	req.URL.RawQuery = q.Encode()

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to authenticate: %s", resp.Status)
// 	}

// 	token := &Token{}

// 	if err = json.NewDecoder(resp.Body).Decode(token); err != nil {
// 		return err
// 	}

// 	c.SetBearerAuth(token.AccessToken)
// 	endpoint := fmt.Sprintf("https://%s/v2/", host)

// 	req, err = http.NewRequest("GET", endpoint, nil)
// 	if err != nil {
// 		return err
// 	}
// 	if c.Credentials != nil {
// 		req.Header.Set("Authorization", c.Credentials.encoded)
// 	}

// 	client = &http.Client{}
// 	resp, err = client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to authenticate: %s", resp.Status)
// 	}

// 	fmt.Println("-- DEBUG: Successfully authenticated with realm: ", token.AccessToken)

// 	return nil
// }

// // GetCredentials retrieves the credentials for the given registry from the docker config
// // and sets the basic auth header for the client. This should make integration into existing CI/CD workflows easier.
// // We may need support for allowing the user to set the credentials manually in the future.
// func (c *Client) GetCredentials(ref string) error {
// 	ctx := context.Background()
// 	reference, err := ParseRef(ref)
// 	if err != nil {
// 		return err
// 	}

// 	credOpts := credentials.StoreOptions{}
// 	store, err := credentials.NewStoreFromDocker(credOpts)
// 	if err != nil {
// 		return err
// 	}

// 	creds, err := store.Get(ctx, reference.Host)
// 	if err != nil {
// 		return err
// 	}
// 	if creds.Password != "" {
// 		c.NewAuth(creds.Username, creds.Password, reference.Host)
// 	} else {
// 		return errors.New("no credentials found for registry")
// 	}

// 	return nil
// }
