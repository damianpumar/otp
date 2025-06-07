package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

func main() {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MiAppDemo",
		AccountName: "user@example.com",
	})
	if err != nil {
		log.Fatalf("Error generando clave: %v", err)
	}

	fmt.Println("Clave secreta base32:", key.Secret())
	otpURL := key.URL()

	// Codificar la URL para usarla en el QR generator
	qrURL := "https://api.qrserver.com/v1/create-qr-code/?data=" + url.QueryEscape(otpURL)
	fmt.Println("Escanea este código QR con Google Authenticator o tu app TOTP favorita:")
	fmt.Println(qrURL)

	go func() {
		for {
			currentCode, err := totp.GenerateCode(key.Secret(), time.Now())
			if err != nil {
				log.Fatalf("Error generando código actual: %v", err)
			}
			fmt.Println("Código actual generado por la app:", currentCode)

			time.Sleep(30 * time.Second)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Ingresa el código OTP: ")
		code, _ := reader.ReadString('\n')
		code = strings.TrimSpace(code)

		if totp.Validate(code, key.Secret()) {
			fmt.Println("Código válido ✅")

		} else {
			fmt.Println("Código inválido ❌, intenta otra vez")
		}
	}
}
