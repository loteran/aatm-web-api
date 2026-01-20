# AATM Web API - Amazing Automatic Torrent Maker

Interface web pour creer des torrents avec integration qBittorrent.

> **Note**: Ce projet est base sur [zedeska/aatm](https://github.com/zedeska/aatm) (application Wails desktop).
> Cette version est une reimplementation en tant qu'API web containerisee avec Docker.

## Fonctionnalites

- Explorateur de fichiers avec navigation complete
- Affichage MediaInfo des fichiers video
- Creation de fichiers .torrent
- Generation de fichiers NFO
- Upload automatique vers qBittorrent
- Upload vers La-Cale (tracker prive)
- Historique des fichiers traites

## Installation

### Avec Docker Compose (recommande)

1. Cloner le repository:
```bash
git clone https://github.com/loteran/aatm-web-api.git
cd aatm-web-api
```

2. Copier et configurer le fichier d'environnement:
```bash
cp .env.example .env
nano .env
```

3. Lancer le container:
```bash
docker compose up -d
```

### Avec Docker Hub

```bash
docker run -d \
  --name aatm \
  -p 8085:8080 \
  -p 8086:8081 \
  -v /:/host:ro \
  -v /your/media/path:/media \
  -v aatm-data:/data \
  -v aatm-qbt:/config/qBittorrent \
  loteran/aatm-web-api:latest
```

## Configuration

### Variables d'environnement (.env)

| Variable | Description | Default |
|----------|-------------|---------|
| `MEDIA_PATH` | Chemin vers vos medias (avec acces ecriture) | `/` |
| `AATM_API_PORT` | Port de l'interface web | `8085` |
| `AATM_QBIT_PORT` | Port du WebUI qBittorrent | `8086` |
| `TZ` | Timezone | `Europe/Paris` |

### Acces

| Service | URL |
|---------|-----|
| Interface AATM | http://localhost:8085 |
| qBittorrent WebUI | http://localhost:8086 |

### Credentials qBittorrent par defaut

- **Username**: `admin`
- **Password**: `adminadmin`

## Structure des volumes

| Volume | Description |
|--------|-------------|
| `/host` | Systeme de fichiers hote (lecture seule) |
| `/media` | Chemin media avec acces ecriture |
| `/data` | Base de donnees et settings |
| `/config/qBittorrent` | Configuration qBittorrent |
| `/torrents` | Fichiers .torrent generes |

## Captures d'ecran

L'interface propose:
- Un explorateur de fichiers pour naviguer dans vos medias
- Un panneau de details avec MediaInfo
- Un workflow de creation de torrent en 5 etapes
- Une page de parametres
- Un historique des fichiers traites

## Credits

Ce projet est base sur [zedeska/aatm](https://github.com/zedeska/aatm), une application desktop Wails.
Cette version reimplemente les fonctionnalites en tant qu'API web containerisee avec Docker.

## Licence

MIT
