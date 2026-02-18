# 🧲 AATM Web API

**AATM Web API** (Amazing Automatic Torrent Maker) est un conteneur Docker avec **interface web** pour créer des fichiers **.torrent** avec **qBittorrent intégré**.

Il permet de naviguer dans vos fichiers, générer des torrents et NFO, et uploader directement vers votre client torrent ou La-Cale.

> 🙏 **Basé sur** [zedeska/aatm](https://github.com/zedeska/aatm) - Merci pour le code original !

---

## ✨ Fonctionnalités

- 🌐 **Interface web** moderne dark mode
- 📁 **Explorateur de fichiers** avec navigation complète
- 🎬 Affichage **MediaInfo** des fichiers vidéo
- 🧲 Création de fichiers `.torrent` (avec progression en temps réel)
- 📝 Génération de fichiers **NFO**
- 📂 **Répertoire de sortie configurable** avec organisation automatique :
  - `{outputDir}/Films/{nom}/` pour les films
  - `{outputDir}/Séries/{nom}/` pour les séries
  - `{outputDir}/Ebooks/{nom}/` pour les ebooks
  - `{outputDir}/Jeux/{nom}/` pour les jeux
- ⬆️ Upload automatique vers **qBittorrent**, **Transmission** ou **Deluge**
- 🚀 Upload vers **La-Cale** (tracker privé) avec aperçu des tags/catégories
- 🔗 Création de **hardlinks** automatique
- ⚙️ Configuration complète via interface web
- 📜 Historique des fichiers traités
- 🐳 qBittorrent inclus dans le conteneur

---

## ⚙️ Variables d'environnement

| Variable | Description | Défaut |
|----------|-------------|--------|
| `MEDIA_PATH` | Chemin vers vos médias sur l'hôte | `/` |
| `AATM_API_PORT` | Port de l'interface web | `8085` |
| `AATM_QBIT_PORT` | Port du WebUI qBittorrent | `8086` |
| `TZ` | Timezone | `Europe/Paris` |

---

## 📁 Volumes

| Chemin conteneur | Description |
|------------------|-------------|
| `/host` | Système de fichiers hôte (lecture seule) |
| `/host/mnt` | `/mnt` hôte avec accès écriture |
| `/host/media` | `/media` hôte avec accès écriture |
| `/host/home` | `/home` hôte avec accès écriture |
| `/data` | Base de données et settings |
| `/config/qBittorrent` | Configuration qBittorrent |
| `/torrents` | Répertoire de sortie par défaut pour les .torrent et .nfo |

---

## 📂 Organisation des fichiers de sortie

Par défaut, les fichiers `.torrent` et `.nfo` sont créés dans `/torrents` (mappé sur `./torrents/` côté hôte), organisés automatiquement :

```
/torrents/
├── Films/
│   └── The.film.2024.MULTi.1080p.BluRay.AC3.5.1.X265-GROUPE/
│       ├── The.film.2024.MULTi.1080p.BluRay.AC3.5.1.X265-GROUPE.torrent
│       └── The.film.2024.MULTi.1080p.BluRay.AC3.5.1.X265-GROUPE.nfo
├── Séries/
├── Ebooks/
└── Jeux/
```

Le répertoire de sortie est configurable dans **Paramètres > Chemins > Répertoire de sortie**.
Vous pouvez pointer vers n'importe quel chemin accessible depuis le conteneur (ex: `/host/mnt/Stockage/Torrents`).

---

## 🔗 Hardlinks - Configuration importante

Pour créer des **hardlinks** (liens physiques) vers vos fichiers, le conteneur a besoin d'un accès en **écriture** aux répertoires concernés.

Par défaut, le système de fichiers est monté en **lecture seule** (`/:/host:ro`) pour la sécurité. Vous devez donc monter explicitement les répertoires où vous souhaitez créer des hardlinks.

### Répertoires courants

Les répertoires suivants sont généralement utilisés pour les médias :
- `/mnt` - Disques montés
- `/media` - Médias
- `/home` - Dossiers utilisateurs

### Ajouter vos propres répertoires

Si vos médias ou répertoires de hardlinks sont ailleurs (ex: `/data`, `/srv`), ajoutez une ligne dans les volumes :

```yaml
volumes:
  - /data:/host/data
  - /srv:/host/srv
```

> ⚠️ **Note** : Les hardlinks ne fonctionnent qu'entre fichiers sur le **même système de fichiers** (même partition/disque).

---

## 🚀 Exemple docker-compose

```yaml
services:
  aatm-web-api:
    image: loteran/aatm-web-api:latest
    container_name: aatm-web-api
    restart: unless-stopped
    ports:
      - "8085:8080"      # Interface web
      - "8086:8081"      # qBittorrent WebUI
      - "6881:6881"      # Torrent port
      - "6881:6881/udp"
    environment:
      - TZ=Europe/Paris
    volumes:
      - ./data:/data
      - ./qbt-config:/config/qBittorrent
      - ./torrents:/torrents
      # Lecture seule pour la navigation
      - /:/host:ro
      # Écriture pour hardlinks/torrents/nfo (ajoutez vos répertoires ici)
      - /mnt:/host/mnt
      - /media:/host/media
      - /home:/host/home
      # Exemple: si vos médias sont dans /data ou /srv
      # - /data:/host/data
      # - /srv:/host/srv
```

---

## 🎛️ Clients torrent supportés

| Client | Support |
|--------|---------|
| qBittorrent | ✅ (intégré dans le conteneur) |
| Transmission | ✅ (instance externe) |
| Deluge | ✅ (instance externe) |
| Aucun | ✅ (désactiver l'upload automatique) |

Configurez le client dans **Paramètres > Client Torrent**.

---

## 🔑 Clé API La-Cale

Pour pouvoir uploader vos torrents sur **La-Cale**, vous devez générer une clé API depuis votre compte :

1. Rendez-vous sur **https://la-cale.space/settings/api-keys**
2. Générez une nouvelle clé API
3. Copiez la clé et renseignez-la dans **Paramètres > La-Cale > Clé API** de l'interface AATM

> ⚠️ **Sans cette clé API, l'upload vers La-Cale ne fonctionnera pas.**

En cas d'échec de l'upload La-Cale, les fichiers `.torrent` et `.nfo` locaux sont conservés et leurs chemins sont affichés. Un bouton **"Terminer sans upload La-Cale"** permet de clôturer le workflow.

---

## 🖥️ Utilisation

1. Lancez le conteneur
2. Accédez à `http://votre-ip:8085`
3. Configurez vos paramètres (client torrent, clé API La-Cale, répertoire de sortie...)
4. Naviguez dans `/host` pour trouver vos fichiers
5. Sélectionnez un fichier vidéo et suivez le workflow :
   - **Étape 1** : Sélection du fichier
   - **Étape 2** : Recherche des métadonnées (TMDB, etc.)
   - **Étape 3** : Génération du NFO
   - **Étape 4** : Création du torrent
   - **Étape 5** : Aperçu et upload vers La-Cale

---

## 🔐 Credentials qBittorrent par défaut

| Paramètre | Valeur |
|-----------|--------|
| URL | `http://localhost:8086` |
| Username | `admin` |
| Password | `adminadmin` |

---

## 📝 Notes

- La configuration est persistante dans `/data`
- qBittorrent est intégré dans le conteneur
- Compatible architectures `amd64` (PC/UNRAID) et `arm64` (Raspberry Pi)

---

## 📋 Changelog

### v4.0.1
- Répertoire de sortie configurable pour les fichiers `.torrent` et `.nfo`
- Organisation automatique en sous-dossiers par type (Films/, Séries/, Ebooks/, Jeux/)
- Chaque release isolée dans son propre sous-répertoire
- Correction : les résultats Transmission et La-Cale s'affichent dans des blocs séparés
- Correction : le workflow ne navigue plus vers la page fichiers en cas d'échec upload La-Cale
- Bouton "Terminer sans upload La-Cale" en cas d'échec

### v4.0.0
- API La-Cale avec aperçu des catégories et tags
- Support Transmission et Deluge en plus de qBittorrent
- Amélioration du workflow en 5 étapes

---

## 🔗 Liens

- **DockerHub** : https://hub.docker.com/repository/docker/loteran/aatm-web-api/general
- **Basé sur** : https://github.com/zedeska/aatm
