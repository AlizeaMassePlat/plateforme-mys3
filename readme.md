# My S3 Clone

My S3 Clone est une API légère qui reproduit les fonctionnalités de base d'un service de stockage de type S3 en utilisant MinIO comme backend de stockage. Elle permet de créer des buckets, télécharger et récupérer des fichiers, lister les fichiers présents dans un bucket et supprimer des fichiers.

## Fonctionnalités

- **Créer un Bucket** : Crée un bucket de stockage dans MinIO.
- **Uploader un Objet** : Télécharge un objet dans un bucket.
- **Lister les Buckets** : Récupère la liste de tous les buckets.
- **Récupérer un Objet** : Récupère un objet spécifique depuis un bucket.
- **Supprimer un Objet** : Supprime un objet d'un bucket.
- **Supprimer un Bucket** : Supprime un bucket de MinIO.

## Prérequis

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Installation et Lancement

1. Clonez le dépôt :

    ```bash
    git clone https://github.com/votre-utilisateur/my-s3-clone.git
    cd my-s3-clone
    ```

2. Construisez et démarrez les conteneurs avec Docker Compose :

    ```bash
    docker-compose up --build
    ```

3. L'API sera accessible sur `http://localhost:8080`.

## Utilisation

Vous pouvez utiliser [Postman](https://www.postman.com/) pour interagir avec l'API. Voici quelques exemples de requêtes :


# Configuration des Requêtes Postman

## 1. Créer un Bucket

- **Type de requête**: `POST`
- **URL**: `http://localhost:8080/buckets`

**Configuration dans Postman**:

1. **Sélectionner la méthode HTTP**: Cliquez sur le menu déroulant à côté de l'URL et choisissez `POST`.
2. **Entrer l'URL**: Saisissez `http://localhost:8080/buckets` dans le champ URL.
3. **Configurer le corps de la requête**:
   - Cliquez sur l'onglet `Body`.
   - Sélectionnez `raw`.
   - Choisissez `XML` dans le menu déroulant situé à droite de la sélection `raw`.
   - Entrez le contenu suivant dans le champ texte :
     ```xml
     <Bucket><Name>myBucket</Name></Bucket>
     ```
4. **Envoyer la requête**: Cliquez sur le bouton `Send`.

## 2. Uploader un Objet dans un Bucket

- **Type de requête**: `PUT`
- **URL**: `http://localhost:8080/buckets/myBucket?object=myObject`

**Configuration dans Postman**:

1. **Sélectionner la méthode HTTP**: Cliquez sur le menu déroulant à côté de l'URL et choisissez `PUT`.
2. **Entrer l'URL**: Saisissez `http://localhost:8080/buckets/myBucket?object=myObject` dans le champ URL.
3. **Configurer le corps de la requête**:
   - Cliquez sur l'onglet `Body`.
   - Sélectionnez `raw`.
   - Choisissez `Text` dans le menu déroulant situé à droite de la sélection `raw`.
   - Entrez le contenu du fichier que vous souhaitez uploader, par exemple :
     ```
     Ceci est le contenu de l'objet.
     ```
4. **Envoyer la requête**: Cliquez sur le bouton `Send`.

## 3. Lister les Buckets

- **Type de requête**: `GET`
- **URL**: `http://localhost:8080/buckets`

**Configuration dans Postman**:

1. **Sélectionner la méthode HTTP**: Cliquez sur le menu déroulant à côté de l'URL et choisissez `GET`.
2. **Entrer l'URL**: Saisissez `http://localhost:8080/buckets` dans le champ URL.
3. **Envoyer la requête**: Cliquez sur le bouton `Send`.

## 4. Télécharger un Objet Spécifique

- **Type de requête**: `GET`
- **URL**: `http://localhost:8080/buckets/myBucket?object=myObject`

**Configuration dans Postman**:

1. **Sélectionner la méthode HTTP**: Cliquez sur le menu déroulant à côté de l'URL et choisissez `GET`.
2. **Entrer l'URL**: Saisissez `http://localhost:8080/buckets/myBucket?object=myObject` dans le champ URL.
3. **Envoyer la requête**: Cliquez sur le bouton `Send`.

## 5. Supprimer un Objet dans un Bucket

- **Type de requête**: `DELETE`
- **URL**: `http://localhost:8080/buckets/myBucket?object=myObject`

**Configuration dans Postman**:

1. **Sélectionner la méthode HTTP**: Cliquez sur le menu déroulant à côté de l'URL et choisissez `DELETE`.
2. **Entrer l'URL**: Saisissez `http://localhost:8080/buckets/myBucket?object=myObject` dans le champ URL.
3. **Envoyer la requête**: Cliquez sur le bouton `Send`.

## 6. Supprimer un Bucket

- **Type de requête**: `DELETE`
- **URL**: `http://localhost:8080/buckets/myBucket`

**Configuration dans Postman**:

1. **Sélectionner la méthode HTTP**: Cliquez sur le menu déroulant à côté de l'URL et choisissez `DELETE`.
2. **Entrer l'URL**: Saisissez `http://localhost:8080/buckets/myBucket` dans le champ URL.
3. **Envoyer la requête**: Cliquez sur le bouton `Send`.

## Résultats Attendus

- **Créer un Bucket**: La réponse doit être HTTP 201 Created.
- **Uploader un Objet**: La réponse doit être HTTP 201 Created.
- **Lister les Buckets**: La réponse doit afficher la liste des buckets en XML.
- **Télécharger un Objet**: La réponse doit contenir le contenu de l'objet.
- **Supprimer un Objet**: La réponse doit être HTTP 204 No Content.
- **Supprimer un Bucket**: La réponse doit être HTTP 204 No Content.


### Visualisation dans l'Interface Graphique

L'interface graphique de MinIO est une couche de visualisation qui permet de voir et de gérer les ressources (buckets et objets) stockées dans MinIO. Les actions effectuées via l'API, telles que la création d'un bucket ou le téléchargement d'un objet, sont immédiatement reflétées dans l'interface graphique.

#### Exemple d'Interaction :

- **Création d'un Bucket via l'API**:
  - Lorsque vous créez un bucket en envoyant une requête `POST` à l'API (`http://localhost:8080/buckets`), MinIO crée le bucket en interne.
  - Si vous ouvrez ensuite l'interface graphique de MinIO (disponible à `http://localhost:9000`), vous verrez le nouveau bucket apparaître dans la liste des buckets.

- **Téléchargement d'un Objet via l'API**:
  - De même, lorsque vous téléchargez un objet dans un bucket en utilisant une requête `PUT` via l'API, cet objet devient immédiatement visible dans le bucket correspondant dans l'interface graphique de MinIO.

- **Suppression via l'Interface Graphique**:
  - Si vous supprimez un bucket ou un objet en utilisant l'interface graphique de MinIO, ces changements seront également reflétés lorsque vous listerez les buckets ou les objets via l'API.

---

### Comment Fonctionne l'Authentification AWS Signature Version 4 ?

L'authentification via AWS Signature Version 4 (SigV4) implique plusieurs étapes qui incluent la création d'une signature cryptographique basée sur la clé secrète de l'utilisateur et les détails de la requête. Voici un résumé des étapes de la signature SigV4, comme implémenté dans le code :

1. **Construction de la Requête Canonique** :
   - Une requête canonique est une version normalisée de la requête HTTP. Elle inclut les informations comme le chemin URI, les en-têtes, les paramètres de requête, et un hachage du corps de la requête. Dans la fonction `createCanonicalRequest`, ces composants sont formatés dans un ordre spécifique.

2. **Création de la Chaîne à Signer** :
   - La "chaîne à signer" est une combinaison de l'algorithme de signature, de l'horodatage de la requête (`x-amz-date`), de la portée des crédences (qui inclut la date et la région), et du hachage SHA-256 de la requête canonique. Cette étape est réalisée dans la fonction `createStringToSign`.

3. **Calcul de la Signature** :
   - La signature est générée en utilisant HMAC-SHA256 pour hacher la "chaîne à signer" avec plusieurs couches de clés dérivées de la clé secrète AWS. Ces couches incluent des informations spécifiques à la date et à la région, ce qui garantit que la signature n'est valide que pour une période et une région spécifiques. Ceci est implémenté dans `calculateSignature`.

4. **Ajout de l'En-tête d'Autorisation** :
   - Enfin, la signature est ajoutée à la requête HTTP dans l'en-tête `Authorization`, formant un en-tête qui inclut les informations d'identification AWS, la portée des crédences, les en-têtes signés, et la signature elle-même. Cela est fait dans la fonction `addS3Authorization`.

### Les En-têtes Obligatoires pour la Signature AWS Version 4

Pour qu'une requête signée soit acceptée par AWS S3, les en-têtes suivants sont obligatoires :

- **Host** : Indique l'endpoint de l'API S3, comme `s3.amazonaws.com` ou une version spécifique à une région, comme `s3.us-east-1.amazonaws.com`.

- **x-amz-date** : Spécifie l'horodatage de la requête au format ISO8601. Par exemple, `20230904T123600Z`.

- **x-amz-content-sha256** : Représente le hachage SHA-256 du corps de la requête. Pour les requêtes sans corps ou où la charge utile n'est pas signée explicitement, utilisez `"UNSIGNED-PAYLOAD"`.

- **Authorization** : Contient les informations de signature générées via AWS SigV4. Cet en-tête prouve que la requête est authentique et autorisée par AWS.

---