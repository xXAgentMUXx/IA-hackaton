package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
)

type UsageData struct {
	DeviceType       string  `json:"deviceType"`
	PhoneModel       string  `json:"phoneModel"`
	Streaming        float64 `json:"streaming"`
	Emails           int     `json:"emails"`
	VideoCalls       float64 `json:"videoCalls"`
	CloudStorage     float64 `json:"cloudStorage"`
	SearchQueries    int     `json:"searchQueries"`
	SocialMediaHours float64 `json:"socialMediaHours"`
	Downloads        float64 `json:"downloads"`
	MusicStreaming   float64 `json:"musicStreaming"`
	PhotoSharing     int     `json:"photoSharing"`
	GPSUsage         float64 `json:"gpsUsage"`
}

type Result struct {
	CO2  float64  `json:"co2"`
	Tips []string `json:"tips"`
}

type Phone struct {
	Model string  `json:"model"`
	CO2   float64 `json:"co2"`
}

var phones []Phone

func loadPhones() {
	data, err := ioutil.ReadFile("phones.json")
	if err != nil {
		log.Fatalf("Erreur lecture phones.json: %v", err)
	}
	if err := json.Unmarshal(data, &phones); err != nil {
		log.Fatalf("Erreur parsing phones.json: %v", err)
	}
}

func getPhoneCO2(model string) (float64, bool) {
	for _, phone := range phones {
		if strings.EqualFold(phone.Model, model) {
			return phone.CO2, true
		}
	}
	return 0, false
}

func calculateCO2(data UsageData) Result {
	co2 := 0.0
	tips := []string{}

	co2 += data.Streaming * 55
	co2 += float64(data.Emails) * 4
	co2 += data.VideoCalls * 50
	co2 += data.CloudStorage * 10
	co2 += float64(data.SearchQueries) * 0.3
	co2 += data.SocialMediaHours * 30
	co2 += data.Downloads * 5
	co2 += data.MusicStreaming * 20
	co2 += float64(data.PhotoSharing) * 1.5
	co2 += data.GPSUsage * 8

	if data.DeviceType == "telephone" {
		if val, ok := getPhoneCO2(data.PhoneModel); ok {
			co2 += val
		} else {
			tips = append(tips, "⚠️ Téléphone non reconnu : modèle générique estimé.")
			co2 += 30 
		}
	}

	if data.Streaming >= 1 {
		tips = append(tips, "Réduisez le streaming (baissez la qualité ou téléchargez).")
	}
	if data.Emails > 10 {
		tips = append(tips, "Supprimez vos anciens e-mails.")
	}
	if data.CloudStorage > 5 {
		tips = append(tips, "Nettoyez votre cloud.")
	}
	if data.VideoCalls > 2 {
		tips = append(tips, "Désactivez la vidéo quand inutile.")
	}
	if data.SearchQueries > 20 {
		tips = append(tips, "Utilisez vos favoris pour moins de recherches.")
	}
	if data.SocialMediaHours > 1 {
		tips = append(tips, "Réduisez votre temps sur les réseaux sociaux.")
	}
	if data.Downloads > 5 {
		tips = append(tips, "Évitez les téléchargements inutiles.")
	}
	if data.PhotoSharing > 20 {
		tips = append(tips, "Compressez les photos avant de les envoyer.")
	}
	if data.GPSUsage > 1 {
		tips = append(tips, "Fermez votre GPS quand vous ne l'utilisez pas.")
	}

	return Result{CO2: co2, Tips: tips}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var usage UsageData
	if err := json.NewDecoder(r.Body).Decode(&usage); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	result := calculateCO2(usage)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func phonesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(phones)
}

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func main() {
	loadPhones()

	mux := http.NewServeMux()
	mux.HandleFunc("/", AccueilHandler)
	mux.HandleFunc("/api/calculate", calculateHandler)
	mux.HandleFunc("/api/phones", phonesHandler)
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	handler := cors.AllowAll().Handler(mux)
	log.Println("Serveur disponible sur http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
