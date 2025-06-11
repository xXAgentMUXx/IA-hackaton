package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/rs/cors"
)

type UsageData struct {
	Streaming        float64 `json:"streaming"`
	Emails           int     `json:"emails"`
	VideoCalls       float64 `json:"videoCalls"`
	CloudStorage     float64 `json:"cloudStorage"`
	SearchQueries    int     `json:"searchQueries"`
	SocialMediaHours float64 `json:"socialMediaHours"`
	Downloads        float64 `json:"downloads"`
}

type Result struct {
	CO2  float64  `json:"co2"`
	Tips []string `json:"tips"`
}

func calculateCO2(data UsageData) Result {
	co2 := data.Streaming*55 + float64(data.Emails)*4 + data.VideoCalls*50 + data.CloudStorage*10 + float64(data.SearchQueries)*0.3 + data.SocialMediaHours*30 + data.Downloads*5
	tips := []string{}
	if data.Streaming >= 1 {
		tips = append(tips, "Réduisez le streaming en baissant la qualité vidéo ou en téléchargeant les contenus.")
	}
	if data.Emails > 10 {
		tips = append(tips, "Évitez les e-mails inutiles et supprimez les anciens messages.")
	}
	if data.CloudStorage > 5 {
		tips = append(tips, "Nettoyez votre stockage cloud régulièrement.")
	}
	if data.VideoCalls > 2 {
		tips = append(tips, "Désactivez la vidéo quand ce n’est pas nécessaire en réunion.")
	}
	if data.SearchQueries > 20 {
		tips = append(tips, "Évitez les recherches inutiles, utilisez des favoris.")
	}
	if data.SocialMediaHours > 1 {
		tips = append(tips, "Limitez le temps passé sur les réseaux sociaux.")
	}
	if data.Downloads > 5 {
		tips = append(tips, "Évitez les téléchargements répétitifs ou inutiles.")
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
	log.Println("Serveur sur http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}
