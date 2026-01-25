# Mémoire AATM - Container Torrent Uploader

## Description du projet
Container Docker qui permet de :
- Créer des fichiers torrent
- Générer des NFO
- Uploader sur le site **la-cale.space**

## Emplacement
`/home/pi/containers/aatm/`

## Fichiers principaux
- `api/static/index.html` - Interface web principale
- `api/main.go` - API backend Go
- `api/integration.go` - Logique d'upload vers La-Cale
- `docker-compose.yml` - Configuration Docker

## Credentials La-Cale (stockés dans la DB)
- Email: mackenzi93@msn.com
- Password: Mackenzi93
- Passkey: REDACTED_PASSKEY

## Modifications effectuées (Janvier 2026)

### Problème initial
- L'upload de torrent fonctionnait mais la **description était vide** sur La-Cale

### Solution implémentée
Réécriture complète de `/api/static/index.html` pour reproduire les fonctionnalités du code original GitHub (https://github.com/zedeska/aatm)

### Fonctionnalités ajoutées

1. **Intégration TMDB**
   - Clé API: `49d8d37e45764e7c6794ed7dd2d896d4`
   - Recherche de films/séries
   - Récupération poster, synopsis, genres, note

2. **Fonction `parseMediaInfo(nfo)`**
   - Extrait depuis le contenu NFO :
     - Container (MKV, MP4)
     - Résolution (2160p, 1080p, 720p)
     - Codec vidéo (x265, x264, AV1)
     - Format audio (DTS-HD MA, TrueHD, Atmos, etc.)
     - Canaux audio
     - Langues audio et sous-titres
     - HDR (Dolby Vision, HDR10, HDR10+)

3. **Fonction `parseReleaseName(name, nfoContent)`**
   - Extrait depuis le nom de fichier :
     - Titre du média
     - Année
     - Saison/Episode (pour les séries)
     - Source (WEB, BluRay, etc.)
     - Groupe de release

4. **Fonction `generatePresentation(data)`**
   - Génère une description HTML complète avec :
     - Poster TMDB
     - Titre, année, genres
     - Note TMDB
     - Synopsis
     - Tableau des spécifications techniques

5. **Interface utilisateur**
   - Sélecteur de type : Film / Episode / Saison
   - Recherche TMDB avec sélection visuelle
   - Affichage des tags de release (badges colorés)
   - Prévisualisation des informations avant upload

## Commandes utiles

```bash
# Redémarrer le container
cd /home/pi/containers/aatm && docker compose restart

# Voir les logs
cd /home/pi/containers/aatm && docker compose logs -f

# Accéder à l'interface web
# http://<IP_RASPBERRY>:<PORT_AATM>
```

## Code source de référence
GitHub original : https://github.com/zedeska/aatm
- `frontend/src/lib/parser.ts` - Parsers originaux
- `frontend/src/lib/presentation.ts` - Générateur de présentation
- `frontend/src/routes/Process.svelte` - Workflow UI

## Modifications Janvier 2026 (suite)

### Support multi-médias ajouté
- **Films/Séries** : TMDB (The Movie Database)
- **E-books** : Google Books API
- **Jeux vidéo** : Steam API

### Extensions supportées
- **Vidéo** : mkv, mp4, avi, wmv, m4v
- **E-books** : epub, pdf, mobi, azw3, cbz, cbr, djvu
- **Jeux** : iso, nsp, xci, pkg, zip, rar, 7z

### Workflow simplifié
- Étape "Terminer" upload automatiquement sur La-Cale
- Auto-détection du type de média
- Recherche TMDB sans l'année (meilleurs résultats)

### Historique
- Endpoint `/api/processed` pour récupérer l'historique depuis la DB
- Affiche tous les fichiers traités avec date

## Notes
- Le container peut afficher "unhealthy" mais fonctionne quand même
- L'API interne La-Cale utilise des endpoints spécifiques pour l'upload
