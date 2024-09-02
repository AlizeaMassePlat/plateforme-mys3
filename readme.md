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
