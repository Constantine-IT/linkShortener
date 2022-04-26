package handlers

import (
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

//	AuthCookie - middleware проверяющее наличие Cookie - "userid" в Request,
//	- если такая cookie там есть, проверяем её на подлинность,
//	- если такой cookie нет, или она не проходит проверку подлинности, то генерируем новый userid,
//	шифруем его с помощью симметричного алгоритма AES и вставляем в Response.
func (app *Application) AuthCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	слайс длиной 16 байт - переменная в которую будет помещаться информация при шифровке/расшифровке
		authCookie := make([]byte, aes.BlockSize)

		//	секретный ключ симметричного шифрования. Длина ключа - 16 байт.
		secretKey := []byte("sbHYDYWgdakkHHDS")

		//	проверочная фраза длиной 6 байт, в связке с userid (длиной 10 байт) - составит суммарно 16 байт
		nonce := []byte("YANDEX")

		// инициализируем интерфейс симметричного шифрования - cipher.Block
		aesblock, err := aes.NewCipher(secretKey)
		if err != nil {
			app.ErrorLog.Fatal(err)
		}

		//	проверяем на наличие в запросе cookie с "userid"
		if requestUserID, err := r.Cookie("userid"); err == nil {
			// если "userid" задан, то проверяем его на подлинность
			// декодируем userid в тип []byte
			requestUserIDByte, err := hex.DecodeString(requestUserID.Value)
			if err != nil {
				app.ErrorLog.Printf("Auth Cookie decoding: %v\n", err)
			}
			// расшифровываем "userid" из Cookie в переменную authCookie, используя созданный ранее cipher.Block AES
			aesblock.Decrypt(authCookie, requestUserIDByte)
			//	проверяем, не заканчивается ли authCookie символами проверочной фразы - nonce
			//	если ДА, то проверка подлинности пройдена
			if string(authCookie[len(authCookie)-len(nonce):]) == string(nonce) {
				next.ServeHTTP(w, r)
				return
			}
		}

		//	если cookie "userid" отсутствует, или не прошло проверку подлинности, то генерируем новый UserID длиной 10 байт,
		userID, err := generateRandom(10)
		if err != nil {
			app.ErrorLog.Printf("UserId generation error: %v\n", err)
			return
		}
		//	добавляем к UserID проверочную фразу - nonce (длиной 6 байт), получаем slice из 16 байт - размер блока для AES.
		//	всё что больше размера блока AES - алгоритм обрезал бы.
		aesblock.Encrypt(authCookie, append(userID, nonce...))	// зашифровываем (UserID + nonce) в переменную authCookie

		//	вставляем зашифрованный (UserID + nonce) в response в виде cookie со сроком жизни - 1 год.
		http.SetCookie(w, &http.Cookie{
			Name: "userid", Value: hex.EncodeToString(authCookie), Expires: time.Now().AddDate(1, 0, 0),
		})

		//	такой же cookie добавляем и в request, чтобы связать сессию с только что созданным userid
		r.AddCookie(&http.Cookie{
			Name: "userid", Value: hex.EncodeToString(authCookie), Expires: time.Now().AddDate(1, 0, 0),
		})

		next.ServeHTTP(w, r)
	})
}

// generateRandom - генерирует случайную последовательность байт длиной size
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
