# Hand-in 7 Report
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/7)
___

## System design
Keychain generation:
1. Generate RSA key with 2048 bits
1. Create 32 bytes SHA256 hash of the secret key 
and save it at the beginning of the file for future consistency check
1. Derive AES256 32 bytes key and 32 bytes salt from the password 
using [scrypt](https://pkg.go.dev/golang.org/x/crypto/scrypt)
1. Write salt as the next bytes to the file
1. Encrypt the message with AES-256 [GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
with the generated key
1. Return public key encoded in base64

Keychain signing:
1. Read the secret key from the file
    1. Read the 32 bytes of SHA256 of the secret key
    1. Read the 32 bytes of the salt 
    1. Read the rest of the file with the ciphertext
1. Use given password, and the read salt to derive the AES256 key
1. Decrypt the ciphertext with generated key
1. Compute SHA256 of the plaintext and compare it with the hash
included in the file, if they don't match abort
1. Sign the data using RSA with the secret key


## Security measures against bruteforce
1. System checks whether the strong password is used - minimum eight characters, 
at least one uppercase letter, one lowercase letter and one number
1. An artificial slowdown if the passwords don't match - 
sleep 5 seconds if it was not possible to decrypt the keychain, 
thus reducing the speed of attacker

## Testing
Please see [test.go](test.go), we basically tried to generate key,
store it in the `secret.enc` encrypted, then to sign the data and 
verify, whether the signature is correct.