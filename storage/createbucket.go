func createBucket(bucketName string) error {
	// Définir l'endpoint et les clés d'accès
	endpoint := "http://minio:9000" // Remplacez par l'endpoint de votre service S3
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	region := "us-east-1" // Remplacez par votre région

	// Construire l'URL de la requête
	url := fmt.Sprintf("%s/%s", endpoint, bucketName)

	// Créer la requête HTTP PUT pour créer un bucket
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	// Ajouter les en-têtes nécessaires pour l'authentification et la compatibilité S3
	req.Header.Add("Host", bucketName)
	req.Header.Add("Date", getCurrentDateHeader()) // Implémentez une fonction pour obtenir la date au format requis
	req.Header.Add("Authorization", generateAuthorizationHeader(req, accessKeyID, secretAccessKey))

	// Envoyer la requête HTTP
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Lire la réponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create bucket: %s", body)
	}

	return nil
}

// Exemple de fonction pour obtenir l'en-tête de la date (à implémenter)
func getCurrentDateHeader() string {
	// Implémentez cette fonction pour retourner la date au format requis par S3
	return "Sat, 30 Sep 2023 00:00:00 GMT"
}

// Exemple de fonction pour générer l'en-tête d'autorisation (à implémenter)
func generateAuthorizationHeader(req *http.Request, accessKeyID, secretAccessKey string) string {
	// Implémentez cette fonction pour générer l'en-tête d'autorisation basé sur votre clé d'accès et clé secrète
	return "AWS " + accessKeyID + ":" + secretAccessKey
}
