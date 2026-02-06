# 🧲 AATM Web API

**AATM Web API** (Amazing Automatic Torrent Maker) est un conteneur Docker avec **interface web** pour créer des fichiers **.torrent** avec **qBittorrent intégré**.

Il permet de naviguer dans vos fichiers, générer des torrents et NFO, et uploader directement vers qBittorrent ou La-Cale.

> 🙏 **Basé sur** [zedeska/aatm](https://github.com/zedeska/aatm) - Merci pour le code original !

---

## ✨ Fonctionnalités

- 🌐 **Interface web** moderne dark mode
- 📁 **Explorateur de fichiers** avec navigation complète
- 🎬 Affichage **MediaInfo** des fichiers vidéo
- 🧲 Création de fichiers `.torrent`
- 📝 Génération de fichiers **NFO**
- ⬆️ Upload automatique vers **qBittorrent** (intégré)
- 🚀 Upload vers **La-Cale** (tracker privé)
- ⚙️ Configuration via interface web
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
| `/media` | Médias avec accès écriture |
| `/data` | Base de données et settings |
| `/config/qBittorrent` | Configuration qBittorrent |
| `/torrents` | Fichiers .torrent générés |

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

## 🔑 Clé API La-Cale

Pour pouvoir uploader vos torrents sur **La-Cale**, vous devez générer une clé API depuis votre compte :

1. Rendez-vous sur **https://la-cale.space/settings/api-keys**
2. Générez une nouvelle clé API
3. Copiez la clé et renseignez-la dans **Paramètres > La-Cale > Clé API** de l'interface AATM

> ⚠️ **Sans cette clé API, l'upload vers La-Cale ne fonctionnera pas.**

---

## 🖥️ Utilisation

1. Lancez le conteneur
2. Accédez à `http://votre-ip:8085`
3. Configurez votre clé API La-Cale dans les paramètres (voir section ci-dessus)
4. Naviguez dans `/host` pour trouver vos fichiers
5. Sélectionnez un fichier vidéo
6. Suivez le workflow de création de torrent
7. Upload automatique vers qBittorrent et/ou La-Cale

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

## 🔗 Liens

- **HubDocker** : https://hub.docker.com/repository/docker/loteran/aatm-web-api/general
- **Basé sur** : https://github.com/zedeska/aatm
