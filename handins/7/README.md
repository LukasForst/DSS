# Hand-in 7 Report
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/7)
___

## System design
The system generates RSA key of 2048 bits, it uses cryptographically secure pseudorandom number generator
coming from golang package `crypto/rand`.
In order to encrypt the private part of the RSA key, 
the system uses AES-256 encryption in [GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode)
mode.
Moreover, the system requires user to use a reasonably strong password,
which is then used by [scrypt](https://pkg.go.dev/golang.org/x/crypto/scrypt)
to derive the AES key use for the encryption itself.

The combination of used algorithms ensures that the system 
achieves desired security properties:
* confidentiality - the secret key is encrypted by strong algorithm
* integrity - the selected encryption algorithm ensures that if the secrets
file was tampered with, the decryption would fail, thus ensuring the integrity
* non-repudiation - the only person, which is able to sign the data, is
the holder of the encrypted secret key and the password, 
this ensures the non-repudiation of the signature 
(unless both, secrets file and the passwords are stolen by the adversary) 

### Keychain generation:
1. Generate RSA key with 2048 bits
1. Derive AES256 32 bytes key and 32 bytes salt from the password 
using [scrypt](https://pkg.go.dev/golang.org/x/crypto/scrypt)
1. Write salt as the first 32 bytes to the file
1. Encrypt the message with AES-256 GCM
with the generated key
1. Return public key encoded in base64

### Keychain signing:
1. Read the secret key from the file
    1. Read the 32 bytes of the salt 
    1. Read the rest of the file with the ciphertext
1. Use given password, and the read salt to derive the AES256 key
1. Decrypt the ciphertext with generated key
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

## Executing the program
Simply run `make run`, the [test.go](test.go) should start.