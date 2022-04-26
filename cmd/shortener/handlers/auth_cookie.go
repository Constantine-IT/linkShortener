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
		//	слайс длиной aes.BlockSize = 16 байт - переменная, в которую будет помещаться информация при шифровке/расшифровке
		authCookie := make([]byte, aes.BlockSize)

		//	секретный ключ симметричного шифрования. Длина ключа - 16 байт.
		secretKey := []byte("sbHYDYWgdakkHHDS")

		//	проверочная фраза длиной 6 байт, в связке с userid (длиной 10 байт) - составит суммарно 16 байт
		nonce := []byte("YANDEX")

		//	инициализируем интерфейс симметричного шифрования - cipher.Block
		aesblock, err := aes.NewCipher(secretKey)
		if err != nil {
			app.ErrorLog.Fatal(err)
		}

		//	проверяем на наличие в запросе cookie с "userid"
		if requestUserID, err := r.Cookie("userid"); err == nil {
			// если "userid" задан, то декодируем "userid" в тип []byte, и проверяем его на подлинность
			requestUserIDByte, err := hex.DecodeString(requestUserID.Value)
			if err != nil {
				app.ErrorLog.Printf("Auth Cookie decoding: %v\n", err)
			}
			// расшифровываем "userid" из Cookie в переменную authCookie, используя созданный ранее cipher.Block AES
			aesblock.Decrypt(authCookie, requestUserIDByte)
			//	проверяем, не заканчивается ли authCookie символами проверочной фразы - nonce
			if string(authCookie[len(authCookie)-len(nonce):]) == string(nonce) {
				next.ServeHTTP(w, r)
				return //	если ДА, то проверка подлинности пройдена
			}
		}
		//	если cookie "userid" отсутствует, или не прошло проверку подлинности, то генерируем новый UserID длиной 10 байт,
		userID, _ := generateRandom(10)
		//	добавляем к нему проверочную фразу - nonce (длиной 6 байт), получаем slice из 16 байт - размер блока для AES.
		//	всё что больше размера блока AES - алгоритм обрезал бы.
		aesblock.Encrypt(authCookie, append(userID, nonce...)) // зашифровываем (UserID + nonce) в переменную authCookie

		//	изготавливаем cookie "userid" с зашифрованным (UserID + nonce), со сроком жизни - 1 год
		cookie := &http.Cookie{
			Name: "userid", Value: hex.EncodeToString(authCookie), Expires: time.Now().AddDate(1, 0, 0),
		}
		//	вставляем cookie в response и в request
		http.SetCookie(w, cookie)
		r.AddCookie(cookie)
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
