package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
)

func main() {
	// 1. TOTP key ကို generate လုပ်ပါ
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MyApp",           // အသုံးပြုမည့် app name
		AccountName: "user@domain.com", // သင့် account name
	})
	if err != nil {
		log.Fatal(err)
	}

	// 2. Generated TOTP Key ကို print လုပ်ပါ
	fmt.Println("Key URL:", key.URL())

	// 3. TOTP token ထုတ်ပါ
	// 30-seconds interval ဖြစ်တဲ့ OTP ကို generate လုပ်ပါ
	otp, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// 4. Generated OTP ကို print လုပ်ပါ
	fmt.Println("Generated OTP:", otp)

	// 5. TOTP Validation (Verification)
	valid := totp.Validate(otp, key.Secret())
	if valid {
		fmt.Println("OTP is valid!")
	} else {
		fmt.Println("OTP is invalid!")
	}
}
