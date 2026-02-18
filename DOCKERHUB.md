# 🧲 AATM Web API

**AATM Web API** (Amazing Automatic Torrent Maker) est un conteneur Docker avec **interface web** pour créer des fichiers **.torrent** avec **qBittorrent intégré**.

Il permet de naviguer dans vos fichiers, générer des torrents et NFO, et uploader directement vers votre client torrent ou La-Cale.

> 🙏 **Basé sur** [zedeska/aatm](https://github.com/zedeska/aatm)

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
| `/data` | Base de données et settings persistants |
| `/config/qBittorrent` | Configuration qBittorrent |
| `/torrents` | Répertoire de sortie par défaut (.torrent et .nfo) |

---

## 📂 Organisation des fichiers de sortie

Les fichiers `.torrent` et `.nfo` sont organisés automatiquement :

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

---

## 🎛️ Clients torrent supportés

| Client | Support |
|--------|---------|
| qBittorrent | ✅ (intégré dans le conteneur) |
| Transmission | ✅ (instance externe) |
| Deluge | ✅ (instance externe) |
| Aucun | ✅ (désactiver l'upload automatique) |

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
      # Écriture pour hardlinks/torrents/nfo
      - /mnt:/host/mnt
      - /media:/host/media
      - /home:/host/home
```

> ℹ️ Pour des médias dans `/data` ou `/srv`, ajoutez : `- /data:/host/data`

---

## 🔑 Clé API La-Cale

1. Rendez-vous sur **https://la-cale.space/settings/api-keys**
2. Générez une nouvelle clé API
3. Renseignez-la dans **Paramètres > La-Cale > Clé API**

En cas d'échec de l'upload, les fichiers locaux sont conservés et un bouton **"Terminer sans upload La-Cale"** est proposé.

---

## 🔐 Credentials qBittorrent par défaut

| Paramètre | Valeur |
|-----------|--------|
| URL | `http://localhost:8086` |
| Username | `admin` |
| Password | `adminadmin` |

---

## 📋 Changelog

### v4.0.1
- Répertoire de sortie configurable pour `.torrent` et `.nfo`
- Organisation automatique en sous-dossiers par type (Films/, Séries/, Ebooks/, Jeux/)
- Correction affichage des statuts Transmission et La-Cale (blocs séparés)
- Correction : pas de redirection en cas d'échec upload La-Cale
- Bouton "Terminer sans upload La-Cale" en cas d'échec

### v4.0.0
- API La-Cale avec aperçu des catégories et tags
- Support Transmission et Deluge
- Workflow en 5 étapes

---

## 🔗 Liens

- **GitHub** : https://github.com/loteran/aatm-web-api
- **Basé sur** : https://github.com/zedeska/aatm
