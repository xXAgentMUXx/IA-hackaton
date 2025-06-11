package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/rs/cors"
)

type UsageData struct {
	DeviceType        string  `json:"deviceType"`
	PhoneModel        string  `json:"phoneModel"`
	Streaming         float64 `json:"streaming"`
	Emails            int     `json:"emails"`
	VideoCalls        float64 `json:"videoCalls"`
	CloudStorage      float64 `json:"cloudStorage"`
	SearchQueries     int     `json:"searchQueries"`
	SocialMediaHours  float64 `json:"socialMediaHours"`
	Downloads         float64 `json:"downloads"`
	MusicStreaming    float64 `json:"musicStreaming"`
	PhotoSharing      int     `json:"photoSharing"`
	GPSUsage          float64 `json:"gpsUsage"`
}

type Result struct {
	CO2  float64  `json:"co2"`
	Tips []string `json:"tips"`
}

func calculateCO2(data UsageData) Result {
	co2 := 0.0
	co2 += data.Streaming * 55
	co2 += float64(data.Emails) * 4
	co2 += data.VideoCalls * 50
	co2 += data.CloudStorage * 10
	co2 += float64(data.SearchQueries) * 0.3
	co2 += data.SocialMediaHours * 30
	co2 += data.Downloads * 5

	co2 += data.MusicStreaming * 20         
	co2 += float64(data.PhotoSharing) * 2   
	co2 += data.GPSUsage * 5                

	
	if data.DeviceType == "telephone" {
		co2 += 15 
		if data.PhoneModel == "fairphone4" {
			co2 -= 5 
		}
	}
	tips := []string{}
	if data.Streaming > 1 {
		tips = append(tips, "Réduisez le streaming vidéo en baissant la qualité ou en téléchargeant les contenus.")
	}
	if data.Emails > 10 {
		tips = append(tips, "Limitez les e-mails inutiles et supprimez les anciens messages.")
	}
	if data.VideoCalls > 2 {
		tips = append(tips, "Coupez la vidéo lors des réunions non essentielles.")
	}
	if data.CloudStorage > 5 {
		tips = append(tips, "Nettoyez votre espace cloud régulièrement.")
	}
	if data.SearchQueries > 20 {
		tips = append(tips, "Utilisez les favoris pour limiter les recherches répétitives.")
	}
	if data.SocialMediaHours > 1 {
		tips = append(tips, "Réduisez votre temps sur les réseaux sociaux.")
	}
	if data.Downloads > 5 {
		tips = append(tips, "Évitez les téléchargements inutiles.")
	}
	if data.MusicStreaming > 2 {
		tips = append(tips, "Téléchargez vos musiques au lieu de les streamer.")
	}
	if data.PhotoSharing > 10 {
		tips = append(tips, "Réduisez le partage excessif de photos.")
	}
	if data.GPSUsage > 1 {
		tips = append(tips, "Désactivez le GPS quand vous ne l’utilisez pas.")
	}

	return Result{
		CO2:  co2,
		Tips: tips,
	}
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

func AccueilHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", calculateHandler)
	mux.HandleFunc("/", AccueilHandler)

	handler := cors.AllowAll().Handler(mux)
	log.Println("✅ Serveur en ligne sur http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
