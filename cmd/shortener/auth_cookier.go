package main

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

// generateRandom - генерирует случайную последовательность байт
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//	AuthCookie - middleware проверяющее наличие Cookie - "userid" в Request,
//	- если она там есть, то вставляем такую же Cookie в Response
//	- если её там нет, то генерируем userid и вставляем его и в Request, и в Response
//	userid при пересылке клиенту шифруется с помощью симметричного алгоритма AES256
func AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	Секретный ключ симметричного шифрования. Длина ключа - 16 байт.
		secretKey := []byte("sbHYDYWgdakkHHDS")
		//	Проверочная фраза длиной 6 байт, в связке с userid (длиной 10 байт) - составит суммарно 16 байт
		nonce := []byte("YANDEX")

		// Инициализируем cipher.Block
		aesblock, err := aes.NewCipher(secretKey)
		if err != nil {
			log.Printf("Creating cipher.Block error: %v\n", err)
			return
		}

		//	Проверяем на наличие в запросе cookie с заданным "userid"
		if requestUserId, err := r.Cookie("userid"); err == nil {
			// если "userid" задан, то проверяем его на подлинность
			// декодируем userid в тип []byte
			requestUserIdByte, err := hex.DecodeString(requestUserId.Value)
			if err != nil {
				log.Printf("error with decoding of request Auth Cookie to []byte: %v\n", err)
			}

			// расшифровываем "userid" в переменную authCookie, используя cipher.Block и secretKey
			authCookie := make([]byte, aes.BlockSize)
			aesblock.Decrypt(authCookie, requestUserIdByte)

			//	проверяем, не заканчивается ли authCookie символами проверочной фразы - nonce
			if string(authCookie[len(authCookie)-len(nonce):]) == string(nonce) {
				//	если так, то присланная cookie - подлинная
				//	выставляем  cookie в response равной cookie "userid" из request
				//	передаём обработку далее в другие handlers
				http.SetCookie(w, requestUserId)
				next.ServeHTTP(w, r)
				return
			}
		}

		//	если cookie "userid" отсутствует, или не прошло проверку подлинности, то
		//	генерируем новый User ID длиной 10 байт, так чтобы вместе с добавленной к нему
		//	проверочной фразой - nonce, мы бы укладывались ровно в 16 байт - размер блока для AES256
		//	так как в данном случае, мы не используем алгоритм GCM, позволяющий шифровать большие массивы информации,
		//	мы ограничены длиной в 16 байт при шифровании.
		userId, err := generateRandom(10)
		if err != nil {
			log.Printf("UserId generation error: %v\n", err)
			return
		}

		// authCookie := hex.EncodeToString(append(userId, nonce...))
		authCookie := make([]byte, aes.BlockSize) // зашифровываем
		aesblock.Encrypt(authCookie, append(userId, nonce...))
		//fmt.Printf("encrypted: %x\n", authCookie)

		// изготавливаем подпись для нашего ключа, используя алгоритм HMAC и функцию SHA256
		//h := hmac.New(sha256.New, secretKey)
		//h.Write(userId)
		//sign := h.Sum(nil)

		http.SetCookie(w, &http.Cookie{
			Name: "userid", Value: hex.EncodeToString(authCookie), Path: "/", Expires: time.Now().AddDate(1, 0, 0),
		})
		
		r.AddCookie(&http.Cookie{
			Name: "userid", Value: hex.EncodeToString(authCookie),
		})

		next.ServeHTTP(w, r)
		return
	})
}
