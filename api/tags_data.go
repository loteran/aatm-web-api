package main

const tagsData = `{
  "quaiprincipalcategories": [
    {
      "name": "Applications",
      "slug": "applications",
      "emplacementsouscategorie": [
        {
          "name": "Logiciels",
          "slug": "software",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Type (Logiciels)",
              "slug": "type-logiciels",
              "tags": [
                {
                  "name": "3D / CAO",
                  "id": ""
                },
                {
                  "name": "Audio / MAO",
                  "id": ""
                },
                {
                  "name": "Bureautique",
                  "id": ""
                },
                {
                  "name": "Développement",
                  "id": ""
                },
                {
                  "name": "Graphisme / Design",
                  "id": ""
                },
                {
                  "name": "Sécurité / Réseau",
                  "id": ""
                },
                {
                  "name": "Utilitaires",
                  "id": ""
                },
                {
                  "name": "Vidéo / Montage",
                  "id": ""
                }
              ]
            },
            {
              "name": "Plateforme (Logiciels)",
              "slug": "plateforme-logiciels",
              "tags": [
                {
                  "name": "Linux",
                  "id": ""
                },
                {
                  "name": "macOS",
                  "id": ""
                },
                {
                  "name": "Windows",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Systèmes",
          "slug": "systemes",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Logiciels)",
              "slug": "plateforme-logiciels",
              "tags": [
                {
                  "name": "BSD / Autres",
                  "id": ""
                },
                {
                  "name": "Linux",
                  "id": ""
                },
                {
                  "name": "macOS",
                  "id": ""
                },
                {
                  "name": "Windows",
                  "id": ""
                }
              ]
            }
          ]
        }
      ],
      "caracteristiques": null
    },
    {
      "name": "Audio",
      "slug": "audio",
      "emplacementsouscategorie": [
        {
          "name": "Audio divers",
          "slug": "audio-divers",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Type (Audio divers)",
              "slug": "type-audio-divers",
              "tags": [
                {
                  "name": "Conférences",
                  "id": ""
                },
                {
                  "name": "Interviews",
                  "id": ""
                },
                {
                  "name": "Livres audio",
                  "id": ""
                },
                {
                  "name": "Podcasts",
                  "id": ""
                }
              ]
            },
            {
              "name": "Type (Musique)",
              "slug": "type-musique",
              "tags": [
                {
                  "name": "BO Animation",
                  "id": ""
                },
                {
                  "name": "BO - Film",
                  "id": ""
                },
                {
                  "name": "BO Jeu Vidéo",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Audiobooks",
          "slug": "audiobooks",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Type (Musique)",
              "slug": "type-musique",
              "tags": [
                {
                  "name": "BO Animation",
                  "id": ""
                },
                {
                  "name": "BO - Film",
                  "id": ""
                },
                {
                  "name": "BO Jeu Vidéo",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Musique",
          "slug": "music",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Type (Musique)",
              "slug": "type-musique",
              "tags": [
                {
                  "name": "Albums",
                  "id": ""
                },
                {
                  "name": "Bande originale",
                  "id": ""
                },
                {
                  "name": "BO Animation",
                  "id": ""
                },
                {
                  "name": "BO - Film",
                  "id": ""
                },
                {
                  "name": "BO Jeu Vidéo",
                  "id": ""
                },
                {
                  "name": "BO / OST",
                  "id": ""
                },
                {
                  "name": "Compilations",
                  "id": ""
                },
                {
                  "name": "Concerts live",
                  "id": ""
                },
                {
                  "name": "Discographies",
                  "id": ""
                },
                {
                  "name": "Singles / EP",
                  "id": ""
                }
              ]
            },
            {
              "name": "Genres (Musique)",
              "slug": "genres-musique",
              "tags": [
                {
                  "name": "Classique",
                  "id": ""
                },
                {
                  "name": "Électro",
                  "id": ""
                },
                {
                  "name": "Jazz / Blues",
                  "id": ""
                },
                {
                  "name": "Pop",
                  "id": ""
                },
                {
                  "name": "Rap / Hip-Hop",
                  "id": ""
                },
                {
                  "name": "Reggae",
                  "id": ""
                },
                {
                  "name": "Rock / Metal",
                  "id": ""
                },
                {
                  "name": "Soul / Funk",
                  "id": ""
                },
                {
                  "name": "World / Folk",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (Audio)",
              "slug": "format-audio",
              "tags": [
                {
                  "name": "AAC (audio)",
                  "id": ""
                },
                {
                  "name": "AIFF",
                  "id": ""
                },
                {
                  "name": "ALAC",
                  "id": ""
                },
                {
                  "name": "FLAC Lossless",
                  "id": ""
                },
                {
                  "name": "MP3 Lossy",
                  "id": ""
                },
                {
                  "name": "OGG",
                  "id": ""
                },
                {
                  "name": "WAV",
                  "id": ""
                }
              ]
            },
            {
              "name": "Qualité / Bitrate",
              "slug": "qualit-bitrate",
              "tags": [
                {
                  "name": "128 kbps",
                  "id": ""
                },
                {
                  "name": "192 kbps",
                  "id": ""
                },
                {
                  "name": "256 kbps",
                  "id": ""
                },
                {
                  "name": "320 kbps",
                  "id": ""
                },
                {
                  "name": "64 kbps",
                  "id": ""
                },
                {
                  "name": "Lossless",
                  "id": ""
                },
                {
                  "name": "V0",
                  "id": ""
                }
              ]
            },
            {
              "name": "Sample rate",
              "slug": "sample-rate",
              "tags": [
                {
                  "name": "176.4 kHz",
                  "id": ""
                },
                {
                  "name": "192 kHz",
                  "id": ""
                },
                {
                  "name": "44.1 kHz",
                  "id": ""
                },
                {
                  "name": "48 kHz",
                  "id": ""
                },
                {
                  "name": "88.2 kHz",
                  "id": ""
                },
                {
                  "name": "96 kHz",
                  "id": ""
                }
              ]
            },
            {
              "name": "Profondeur",
              "slug": "profondeur",
              "tags": [
                {
                  "name": "16-bit",
                  "id": ""
                },
                {
                  "name": "24-bit",
                  "id": ""
                },
                {
                  "name": "32-bit",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Podcast",
          "slug": "podcast",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Type (Musique)",
              "slug": "type-musique",
              "tags": [
                {
                  "name": "BO Animation",
                  "id": ""
                },
                {
                  "name": "BO - Film",
                  "id": ""
                },
                {
                  "name": "BO Jeu Vidéo",
                  "id": ""
                }
              ]
            }
          ]
        }
      ],
      "caracteristiques": null
    },
    {
      "name": "Autres",
      "slug": "autres",
      "emplacementsouscategorie": null,
      "caracteristiques": [
        {
          "name": "Divers",
          "slug": "divers",
          "tags": [
            {
              "name": "Imprimante 3d",
              "id": ""
            },
            {
              "name": "Packs audio",
              "id": ""
            },
            {
              "name": "Ressources graphiques",
              "id": ""
            },
            {
              "name": "Templates",
              "id": ""
            },
            {
              "name": "Wallpapers",
              "id": ""
            }
          ]
        }
      ]
    },
    {
      "name": "E-books",
      "slug": "e-books",
      "emplacementsouscategorie": [
        {
          "name": "BD",
          "slug": "bd",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genre (E-Books)",
              "slug": "genre-e-books",
              "tags": [
                {
                  "name": "Adaptation",
                  "id": ""
                },
                {
                  "name": "Aventure (Livres)",
                  "id": ""
                },
                {
                  "name": "Biographie",
                  "id": ""
                },
                {
                  "name": "Documentaire (Livres)",
                  "id": ""
                },
                {
                  "name": "Drame (Livres)",
                  "id": ""
                },
                {
                  "name": "Fantastique (Livres)",
                  "id": ""
                },
                {
                  "name": "Historique (livres)",
                  "id": ""
                },
                {
                  "name": "Humour",
                  "id": ""
                },
                {
                  "name": "Jeunesse",
                  "id": ""
                },
                {
                  "name": "Policier / Thriller (Livres)",
                  "id": ""
                },
                {
                  "name": "Super-héros",
                  "id": ""
                },
                {
                  "name": "Western (Livres)",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues",
              "slug": "langues",
              "tags": [
                {
                  "name": "Anglais",
                  "id": ""
                },
                {
                  "name": "Autres Langues",
                  "id": ""
                },
                {
                  "name": "Français",
                  "id": ""
                },
                {
                  "name": "Japonais",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (E-books)",
              "slug": "format-e-books",
              "tags": [
                {
                  "name": "Autres Formats",
                  "id": ""
                },
                {
                  "name": "AZW",
                  "id": ""
                },
                {
                  "name": "CB7",
                  "id": ""
                },
                {
                  "name": "CBA",
                  "id": ""
                },
                {
                  "name": "CBR",
                  "id": ""
                },
                {
                  "name": "CBT",
                  "id": ""
                },
                {
                  "name": "CBZ",
                  "id": ""
                },
                {
                  "name": "DOC",
                  "id": ""
                },
                {
                  "name": "EPUB",
                  "id": ""
                },
                {
                  "name": "MOBI",
                  "id": ""
                },
                {
                  "name": "PDF",
                  "id": ""
                },
                {
                  "name": "PRC",
                  "id": ""
                },
                {
                  "name": "TXT",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Comics",
          "slug": "comics",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genre (E-Books)",
              "slug": "genre-e-books",
              "tags": [
                {
                  "name": "Action (Livres)",
                  "id": ""
                },
                {
                  "name": "Aventure (Livres)",
                  "id": ""
                },
                {
                  "name": "Drame (Livres)",
                  "id": ""
                },
                {
                  "name": "Fantastique (Livres)",
                  "id": ""
                },
                {
                  "name": "Historique (livres)",
                  "id": ""
                },
                {
                  "name": "Horreur (Livres)",
                  "id": ""
                },
                {
                  "name": "Humour",
                  "id": ""
                },
                {
                  "name": "Polar (Livres)",
                  "id": ""
                },
                {
                  "name": "Science-fiction (Livres)",
                  "id": ""
                },
                {
                  "name": "Super-héros",
                  "id": ""
                },
                {
                  "name": "Super-vilains",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues",
              "slug": "langues",
              "tags": [
                {
                  "name": "Anglais",
                  "id": ""
                },
                {
                  "name": "Autres Langues",
                  "id": ""
                },
                {
                  "name": "Français",
                  "id": ""
                },
                {
                  "name": "Japonais",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (E-books)",
              "slug": "format-e-books",
              "tags": [
                {
                  "name": "Autres Formats",
                  "id": ""
                },
                {
                  "name": "AZW",
                  "id": ""
                },
                {
                  "name": "CB7",
                  "id": ""
                },
                {
                  "name": "CBA",
                  "id": ""
                },
                {
                  "name": "CBR",
                  "id": ""
                },
                {
                  "name": "CBT",
                  "id": ""
                },
                {
                  "name": "CBZ",
                  "id": ""
                },
                {
                  "name": "DOC",
                  "id": ""
                },
                {
                  "name": "EPUB",
                  "id": ""
                },
                {
                  "name": "MOBI",
                  "id": ""
                },
                {
                  "name": "PDF",
                  "id": ""
                },
                {
                  "name": "PRC",
                  "id": ""
                },
                {
                  "name": "TXT",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Divers",
          "slug": "divers",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genre (E-Books)",
              "slug": "genre-e-books",
              "tags": [
                {
                  "name": "Biographie",
                  "id": ""
                },
                {
                  "name": "Bricolage",
                  "id": ""
                },
                {
                  "name": "Cinéma",
                  "id": ""
                },
                {
                  "name": "Cuisine",
                  "id": ""
                },
                {
                  "name": "Documentaire (Livres)",
                  "id": ""
                },
                {
                  "name": "Économie",
                  "id": ""
                },
                {
                  "name": "Historique (livres)",
                  "id": ""
                },
                {
                  "name": "Horreur (Livres)",
                  "id": ""
                },
                {
                  "name": "Humour",
                  "id": ""
                },
                {
                  "name": "Informatique",
                  "id": ""
                },
                {
                  "name": "Jardinage",
                  "id": ""
                },
                {
                  "name": "Jeux vidéo",
                  "id": ""
                },
                {
                  "name": "Musique",
                  "id": ""
                },
                {
                  "name": "Photographie",
                  "id": ""
                },
                {
                  "name": "Policier / Thriller (Livres)",
                  "id": ""
                },
                {
                  "name": "Politique",
                  "id": ""
                },
                {
                  "name": "Santé \u0026 Bien-être",
                  "id": ""
                },
                {
                  "name": "Sciences",
                  "id": ""
                },
                {
                  "name": "Société",
                  "id": ""
                },
                {
                  "name": "Sport (Livres)",
                  "id": ""
                },
                {
                  "name": "Super-héros",
                  "id": ""
                },
                {
                  "name": "Voyage",
                  "id": ""
                },
                {
                  "name": "Western (Livres)",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues",
              "slug": "langues",
              "tags": [
                {
                  "name": "Anglais",
                  "id": ""
                },
                {
                  "name": "Autres Langues",
                  "id": ""
                },
                {
                  "name": "Français",
                  "id": ""
                },
                {
                  "name": "Japonais",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (E-books)",
              "slug": "format-e-books",
              "tags": [
                {
                  "name": "Autres Formats",
                  "id": ""
                },
                {
                  "name": "AZW",
                  "id": ""
                },
                {
                  "name": "CB7",
                  "id": ""
                },
                {
                  "name": "CBA",
                  "id": ""
                },
                {
                  "name": "CBR",
                  "id": ""
                },
                {
                  "name": "CBT",
                  "id": ""
                },
                {
                  "name": "CBZ",
                  "id": ""
                },
                {
                  "name": "DOC",
                  "id": ""
                },
                {
                  "name": "EPUB",
                  "id": ""
                },
                {
                  "name": "MOBI",
                  "id": ""
                },
                {
                  "name": "PDF",
                  "id": ""
                },
                {
                  "name": "PRC",
                  "id": ""
                },
                {
                  "name": "TXT",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Livres",
          "slug": "livres",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genre (E-Books)",
              "slug": "genre-e-books",
              "tags": [
                {
                  "name": "Action (Livres)",
                  "id": ""
                },
                {
                  "name": "Aventure (Livres)",
                  "id": ""
                },
                {
                  "name": "Biographie",
                  "id": ""
                },
                {
                  "name": "Cuisine",
                  "id": ""
                },
                {
                  "name": "Documentaire (Livres)",
                  "id": ""
                },
                {
                  "name": "Drame (Livres)",
                  "id": ""
                },
                {
                  "name": "Essai",
                  "id": ""
                },
                {
                  "name": "Fantastique (Livres)",
                  "id": ""
                },
                {
                  "name": "Historique (livres)",
                  "id": ""
                },
                {
                  "name": "Horreur (Livres)",
                  "id": ""
                },
                {
                  "name": "Informatique",
                  "id": ""
                },
                {
                  "name": "Jardinage",
                  "id": ""
                },
                {
                  "name": "Jeu de rôle",
                  "id": ""
                },
                {
                  "name": "Jeunesse",
                  "id": ""
                },
                {
                  "name": "Littérature",
                  "id": ""
                },
                {
                  "name": "Philosophie",
                  "id": ""
                },
                {
                  "name": "Polar (Livres)",
                  "id": ""
                },
                {
                  "name": "Policier / Thriller (Livres)",
                  "id": ""
                },
                {
                  "name": "Politique",
                  "id": ""
                },
                {
                  "name": "Psychologique",
                  "id": ""
                },
                {
                  "name": "Romance (Livres)",
                  "id": ""
                },
                {
                  "name": "Romans",
                  "id": ""
                },
                {
                  "name": "Santé \u0026 Bien-être",
                  "id": ""
                },
                {
                  "name": "Science-fiction (Livres)",
                  "id": ""
                },
                {
                  "name": "Sciences",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues",
              "slug": "langues",
              "tags": [
                {
                  "name": "Anglais",
                  "id": ""
                },
                {
                  "name": "Autres Langues",
                  "id": ""
                },
                {
                  "name": "Français",
                  "id": ""
                },
                {
                  "name": "Japonais",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (E-books)",
              "slug": "format-e-books",
              "tags": [
                {
                  "name": "Autres Formats",
                  "id": ""
                },
                {
                  "name": "AZW",
                  "id": ""
                },
                {
                  "name": "CB7",
                  "id": ""
                },
                {
                  "name": "CBA",
                  "id": ""
                },
                {
                  "name": "CBR",
                  "id": ""
                },
                {
                  "name": "CBT",
                  "id": ""
                },
                {
                  "name": "CBZ",
                  "id": ""
                },
                {
                  "name": "DOC",
                  "id": ""
                },
                {
                  "name": "EPUB",
                  "id": ""
                },
                {
                  "name": "MOBI",
                  "id": ""
                },
                {
                  "name": "PDF",
                  "id": ""
                },
                {
                  "name": "PRC",
                  "id": ""
                },
                {
                  "name": "TXT",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Mangas",
          "slug": "mangas",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genre (E-Books)",
              "slug": "genre-e-books",
              "tags": [
                {
                  "name": "Action (Livres)",
                  "id": ""
                },
                {
                  "name": "Aventure (Livres)",
                  "id": ""
                },
                {
                  "name": "Historique (livres)",
                  "id": ""
                },
                {
                  "name": "Horreur (Livres)",
                  "id": ""
                },
                {
                  "name": "Jeunesse",
                  "id": ""
                },
                {
                  "name": "Josei",
                  "id": ""
                },
                {
                  "name": "Kodomo",
                  "id": ""
                },
                {
                  "name": "Psychologique",
                  "id": ""
                },
                {
                  "name": "Romance (Livres)",
                  "id": ""
                },
                {
                  "name": "Science-fiction (Livres)",
                  "id": ""
                },
                {
                  "name": "Seinen",
                  "id": ""
                },
                {
                  "name": "Shōjo",
                  "id": ""
                },
                {
                  "name": "Shōnen",
                  "id": ""
                },
                {
                  "name": "Slice of Life",
                  "id": ""
                },
                {
                  "name": "Sport (Livres)",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues",
              "slug": "langues",
              "tags": [
                {
                  "name": "Anglais",
                  "id": ""
                },
                {
                  "name": "Autres Langues",
                  "id": ""
                },
                {
                  "name": "Français",
                  "id": ""
                },
                {
                  "name": "Japonais",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (E-books)",
              "slug": "format-e-books",
              "tags": [
                {
                  "name": "Autres Formats",
                  "id": ""
                },
                {
                  "name": "AZW",
                  "id": ""
                },
                {
                  "name": "CB7",
                  "id": ""
                },
                {
                  "name": "CBA",
                  "id": ""
                },
                {
                  "name": "CBR",
                  "id": ""
                },
                {
                  "name": "CBT",
                  "id": ""
                },
                {
                  "name": "CBZ",
                  "id": ""
                },
                {
                  "name": "DOC",
                  "id": ""
                },
                {
                  "name": "EPUB",
                  "id": ""
                },
                {
                  "name": "MOBI",
                  "id": ""
                },
                {
                  "name": "PDF",
                  "id": ""
                },
                {
                  "name": "PRC",
                  "id": ""
                },
                {
                  "name": "TXT",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Presse",
          "slug": "presse",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genre (E-Books)",
              "slug": "genre-e-books",
              "tags": [
                {
                  "name": "Actualités",
                  "id": ""
                },
                {
                  "name": "Cinéma",
                  "id": ""
                },
                {
                  "name": "Cuisine",
                  "id": ""
                },
                {
                  "name": "Économie",
                  "id": ""
                },
                {
                  "name": "Historique (livres)",
                  "id": ""
                },
                {
                  "name": "Jardinage",
                  "id": ""
                },
                {
                  "name": "Jeux vidéo",
                  "id": ""
                },
                {
                  "name": "Journaux",
                  "id": ""
                },
                {
                  "name": "Magazines",
                  "id": ""
                },
                {
                  "name": "Politique",
                  "id": ""
                },
                {
                  "name": "Santé \u0026 Bien-être",
                  "id": ""
                },
                {
                  "name": "Sciences",
                  "id": ""
                },
                {
                  "name": "Société",
                  "id": ""
                },
                {
                  "name": "Sport (Livres)",
                  "id": ""
                },
                {
                  "name": "Voyage",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues",
              "slug": "langues",
              "tags": [
                {
                  "name": "Anglais",
                  "id": ""
                },
                {
                  "name": "Autres Langues",
                  "id": ""
                },
                {
                  "name": "Français",
                  "id": ""
                },
                {
                  "name": "Japonais",
                  "id": ""
                }
              ]
            },
            {
              "name": "Format (E-books)",
              "slug": "format-e-books",
              "tags": [
                {
                  "name": "Autres Formats",
                  "id": ""
                },
                {
                  "name": "AZW",
                  "id": ""
                },
                {
                  "name": "CB7",
                  "id": ""
                },
                {
                  "name": "CBA",
                  "id": ""
                },
                {
                  "name": "CBR",
                  "id": ""
                },
                {
                  "name": "CBT",
                  "id": ""
                },
                {
                  "name": "CBZ",
                  "id": ""
                },
                {
                  "name": "DOC",
                  "id": ""
                },
                {
                  "name": "EPUB",
                  "id": ""
                },
                {
                  "name": "MOBI",
                  "id": ""
                },
                {
                  "name": "PDF",
                  "id": ""
                },
                {
                  "name": "PRC",
                  "id": ""
                },
                {
                  "name": "TXT",
                  "id": ""
                }
              ]
            }
          ]
        }
      ],
      "caracteristiques": null
    },
    {
      "name": "Jeux",
      "slug": "games",
      "emplacementsouscategorie": [
        {
          "name": "Nintendo",
          "slug": "nintendo",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Jeux)",
              "slug": "plateforme-jeux",
              "tags": [
                {
                  "name": "3DS",
                  "id": ""
                },
                {
                  "name": "DS/DSi",
                  "id": ""
                },
                {
                  "name": "NES",
                  "id": ""
                },
                {
                  "name": "SNES",
                  "id": ""
                },
                {
                  "name": "Switch",
                  "id": ""
                },
                {
                  "name": "Wii",
                  "id": ""
                },
                {
                  "name": "Wii U",
                  "id": ""
                }
              ]
            },
            {
              "name": "Extension",
              "slug": "extension",
              "tags": [
                {
                  "name": "ISO",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "PC",
          "slug": "pc",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Jeux)",
              "slug": "plateforme-jeux",
              "tags": [
                {
                  "name": "Arcade / MAME",
                  "id": ""
                },
                {
                  "name": "Emulation",
                  "id": ""
                },
                {
                  "name": "Linux (Jeux)",
                  "id": ""
                },
                {
                  "name": "macOS (Jeux)",
                  "id": ""
                },
                {
                  "name": "VR",
                  "id": ""
                },
                {
                  "name": "Windows (Jeux)",
                  "id": ""
                }
              ]
            },
            {
              "name": "Extension",
              "slug": "extension",
              "tags": [
                {
                  "name": "ISO",
                  "id": ""
                },
                {
                  "name": "Pré Activé",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Playstation",
          "slug": "playstation",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Jeux)",
              "slug": "plateforme-jeux",
              "tags": [
                {
                  "name": "PS1",
                  "id": ""
                },
                {
                  "name": "PS2",
                  "id": ""
                },
                {
                  "name": "PS3",
                  "id": ""
                },
                {
                  "name": "PS4",
                  "id": ""
                },
                {
                  "name": "PS5",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Sega",
          "slug": "sega",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Jeux)",
              "slug": "plateforme-jeux",
              "tags": [
                {
                  "name": "Mega Drive",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Téléphones",
          "slug": "t-l-phones",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Jeux)",
              "slug": "plateforme-jeux",
              "tags": [
                {
                  "name": "Android",
                  "id": ""
                },
                {
                  "name": "iOS",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Xbox",
          "slug": "xbox",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Plateforme (Jeux)",
              "slug": "plateforme-jeux",
              "tags": [
                {
                  "name": "Series X|S",
                  "id": ""
                },
                {
                  "name": "Xbox 360",
                  "id": ""
                },
                {
                  "name": "Xbox One",
                  "id": ""
                }
              ]
            }
          ]
        }
      ],
      "caracteristiques": null
    },
    {
      "name": "Vidéo",
      "slug": "video",
      "emplacementsouscategorie": [
        {
          "name": "Films",
          "slug": "films",
          "id": "cmjoyv2cd00027eryreyk39gz",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genres",
              "slug": "genres",
              "tags": [
                {
                  "name": "Action",
                  "id": "cmjudwh76000guyrus1jc9jxs"
                },
                {
                  "name": "Animation",
                  "id": "cmjudwhrd000iuyruhm4bd82d"
                },
                {
                  "name": "Aventure",
                  "id": "cmjudwhi0000huyru6x3xirry"
                },
                {
                  "name": "Biopic",
                  "id": "cmjudwl30000tuyrug7nl1dih"
                },
                {
                  "name": "Collections",
                  "id": "d5l48s1ua6hc738ejoo0"
                },
                {
                  "name": "Comédie",
                  "id": "cmjudwi0v000juyrujjoqy8ya"
                },
                {
                  "name": "Courts-métrages",
                  "id": "cmjudwld7000uuyrusigq41k6"
                },
                {
                  "name": "Documentaire",
                  "id": "cmjudwiay000kuyruuh6y6sux"
                },
                {
                  "name": "Drame",
                  "id": "cmjudwilv000luyruldwrewad"
                },
                {
                  "name": "Émission TV",
                  "id": "cmjudwteq001juyruql1f2o18"
                },
                {
                  "name": "Fantastique",
                  "id": "cmjudwiw7000muyruiwezf9u9"
                },
                {
                  "name": "Guerre",
                  "id": "cmjudwkto000suyruryh2tokp"
                },
                {
                  "name": "Historique",
                  "id": "cmjudwllt000vuyruryo0oj20"
                },
                {
                  "name": "Horreur",
                  "id": "cmjudwj5c000nuyruu32px04l"
                },
                {
                  "name": "Policier / Thriller",
                  "id": "cmjudwjf9000ouyrud0gadn2q"
                },
                {
                  "name": "Romance",
                  "id": "cmjudwkdw000ruyruwu1f7ps7"
                },
                {
                  "name": "Science-fiction",
                  "id": "cmjudwjri000puyru2t0vy9w9"
                },
                {
                  "name": "Sport",
                  "id": "cmjudwt42001iuyru9xod9658"
                },
                {
                  "name": "Téléfilm",
                  "id": "d5hn26pua6hc738iv6ig"
                },
                {
                  "name": "Western",
                  "id": "cmjudwk3e000quyrumxlv4nxg"
                }
              ]
            },
            {
              "name": "Qualité / Résolution",
              "slug": "qualit-r-solution",
              "tags": [
                {
                  "name": "1080p (Full HD)",
                  "id": "44d9f0e4-5ff2-4c6f-9e17-ec2e33acf448"
                },
                {
                  "name": "2160p (4K)",
                  "id": "e42ac76f-d1b0-4ea4-a3da-f332af0f8f4c"
                },
                {
                  "name": "4320p (8K)",
                  "id": "20d145d3-db89-4996-933c-48cf51ee5b6b"
                },
                {
                  "name": "720p (HD)",
                  "id": "86e547c3-7416-46c5-a759-fc1a5d32cc23"
                },
                {
                  "name": "SD",
                  "id": "cmjudwm65000xuyrubkb2ztkf"
                }
              ]
            },
            {
              "name": "Codec vidéo",
              "slug": "codec-vid-o",
              "tags": [
                {
                  "name": "AV1",
                  "id": "cmjudwnvj0010uyrusdqdsydj"
                },
                {
                  "name": "AVC/H264/x264",
                  "id": "cmjoyv2id000u7eryugoe1bee"
                },
                {
                  "name": "HEVC/H265/x265",
                  "id": "cmjoyv2ig000v7eryc9hf1hsa"
                },
                {
                  "name": "MPEG",
                  "id": "5c9bb557-4fd7-488a-8ce7-83c26c7049f2"
                },
                {
                  "name": "VC-1",
                  "id": "168bff0e-8649-4fb4-86d9-87006c2aa511"
                },
                {
                  "name": "VCC/H266/x266",
                  "id": "f77bbe91-82c6-440d-a302-c3015edc19a8"
                },
                {
                  "name": "VP9",
                  "id": "34d10b63-fcd7-4b12-91df-48cde5cb63c0"
                }
              ]
            },
            {
              "name": "Caractéristiques vidéo",
              "slug": "caract-ristiques-vid-o",
              "tags": [
                {
                  "name": "10 bits",
                  "id": "d3827daf-6d53-4d61-8e88-5a5b33ddd85f"
                },
                {
                  "name": "3D",
                  "id": "d5eg09psup7s739bto70"
                },
                {
                  "name": "Dolby Vision",
                  "id": "862f0aaf-e08d-492f-bf2d-9806c74b676a"
                },
                {
                  "name": "HDR",
                  "id": "4e2f5500-05f9-4f35-b273-233abc8ff991"
                },
                {
                  "name": "HDR10+",
                  "id": "da699e63-de34-4c75-8d65-243d8eb51151"
                },
                {
                  "name": "HLG",
                  "id": "79771165-c073-490d-81e3-a408b01cfaa6"
                },
                {
                  "name": "IMAX",
                  "id": "d5gh4ohsup7s73b5irt0"
                },
                {
                  "name": "SDR",
                  "id": "2a3d5386-ba5f-4cce-b957-1456a6a938da"
                }
              ]
            },
            {
              "name": "Source / Type",
              "slug": "source-type",
              "tags": [
                {
                  "name": "4KLight",
                  "id": "d77a8e79-0035-4df8-8cef-8311a1ea1919"
                },
                {
                  "name": "BluRay",
                  "id": "a063935d-adc5-4e69-af04-54f318a80771"
                },
                {
                  "name": "DVDRip",
                  "id": "48c926bc-8fb6-4163-99b7-c745aba4f1ed"
                },
                {
                  "name": "FULL Disc",
                  "id": "fca2b774-3587-426c-8b7a-08e0fe51f2c6"
                },
                {
                  "name": "HDLight",
                  "id": "a1c90ed0-a4b9-4444-bd8d-a4fce8dc5caf"
                },
                {
                  "name": "REMUX",
                  "id": "f4e1b729-bc12-4957-94ee-53d10bfbf63a"
                },
                {
                  "name": "TV",
                  "id": "f8d828b1-1d48-443f-8749-dd4af5cba373"
                },
                {
                  "name": "WEB-DL",
                  "id": "7bd8b291-6e18-4322-9c35-b3470c90039e"
                },
                {
                  "name": "WEBRip",
                  "id": "86413ec8-af63-4fe2-8c88-cb56ec1b176c"
                }
              ]
            },
            {
              "name": "Codec audio",
              "slug": "codec-audio",
              "tags": [
                {
                  "name": "AAC",
                  "id": "cmjudwo5h0011uyruofwyb2cn"
                },
                {
                  "name": "AC3",
                  "id": "cmjudwodw0012uyrut8jh456x"
                },
                {
                  "name": "AC4",
                  "id": "d5e5hc1sup7s73ag6eig"
                },
                {
                  "name": "Autres",
                  "id": "d5e5hnhsup7s73ecfpl0"
                },
                {
                  "name": "DTS",
                  "id": "cmjudwomi0013uyruthukb2pf"
                },
                {
                  "name": "DTS-HD HR",
                  "id": "5e8475fe-db59-4dd8-84fe-7186eaba134d"
                },
                {
                  "name": "DTS-HD MA",
                  "id": "f87c6b8e-6edf-4dbd-b642-205d1060c0d1"
                },
                {
                  "name": "DTS:X",
                  "id": "ec308312-a650-4b00-b011-e624c8779939"
                },
                {
                  "name": "E-AC3",
                  "id": "3ab66a11-d74c-49b7-a3f9-2efc0edfefb2"
                },
                {
                  "name": "E-AC3 Atmos",
                  "id": "38ea7ff5-bd03-4b41-af5c-211792c0d1e8"
                },
                {
                  "name": "FLAC",
                  "id": "d5e6q01sup7s738q0tsg"
                },
                {
                  "name": "HE-AAC",
                  "id": "c1f82902-ad95-4a1c-abdb-3dd8f12b1f40"
                },
                {
                  "name": "MP3",
                  "id": "d5e6q2psup7s738q0tt0"
                },
                {
                  "name": "Opus",
                  "id": "ebad8b2b-0f23-4a3a-a741-9098e61d5f46"
                },
                {
                  "name": "PCM",
                  "id": "4b6f7605-3aa3-411e-8bb4-fd6fcbc9c921"
                },
                {
                  "name": "TrueHD",
                  "id": "c45a5dbb-aad2-461c-96d2-ee4cee3f5d5b"
                },
                {
                  "name": "TrueHD Atmos",
                  "id": "963a0227-d20b-4878-aa84-1deac1de4479"
                }
              ]
            },
            {
              "name": "Langues audio",
              "slug": "langues-audio",
              "tags": [
                {
                  "name": "Autre Langue",
                  "id": "d5elgi9sup7s73emca6g"
                },
                {
                  "name": "Chinois",
                  "id": "d5eldm1sup7s7387o7ng"
                },
                {
                  "name": "English",
                  "id": "d5e5kn1sup7s73b1q4pg"
                },
                {
                  "name": "French",
                  "id": "2d45b5d1-2dfe-4de5-b0fe-e4a08132d06a"
                },
                {
                  "name": "Italian",
                  "id": "d5e5jh1sup7s73ecfppg"
                },
                {
                  "name": "Japanese",
                  "id": "d5e5jchsup7s73ecfpog"
                },
                {
                  "name": "Korean",
                  "id": "d5e5jepsup7s73f6hn9g"
                },
                {
                  "name": "MULTI",
                  "id": "cmjoyv2i0000q7eryz0dnb7f9"
                },
                {
                  "name": "Sans Dialogue",
                  "id": "d5hb34pua6hc739gjva0"
                },
                {
                  "name": "Spanish",
                  "id": "d5e5mhhsup7s73cq0sr0"
                },
                {
                  "name": "VFF",
                  "id": "5cc1e14d-b23a-42c5-956e-3ad9e1529b3c"
                },
                {
                  "name": "VFQ",
                  "id": "b4661137-ed7a-4742-8ebc-40fcb6fa8f8d"
                }
              ]
            },
            {
              "name": "Sous-titres",
              "slug": "sous-titres",
              "tags": [
                {
                  "name": "Autres sous-titres",
                  "id": "d5elg3psup7s73emvf3g"
                },
                {
                  "name": "ENG",
                  "id": "d5e6lu9sup7s73a9mimg"
                },
                {
                  "name": "FR",
                  "id": "cmjudwg9q000euyruscojd2q4"
                },
                {
                  "name": "VFF Sous-Titres",
                  "id": "d5elf41sup7s73emc85g"
                },
                {
                  "name": "VFQ Sous-Titres",
                  "id": "d5elilhsup7s73avoesg"
                }
              ]
            },
            {
              "name": "Extension",
              "slug": "extension",
              "tags": [
                {
                  "name": "Autres Extensions",
                  "id": "d5fuf51sup7s73bbcmq0"
                },
                {
                  "name": "AVI",
                  "id": "d5fueopsup7s73cg5qo0"
                },
                {
                  "name": "ISO",
                  "id": "d5fuevpsup7s739du1l0"
                },
                {
                  "name": "MKV",
                  "id": "d5fuer9sup7s73eq24s0"
                },
                {
                  "name": "MP4",
                  "id": "d5fuelpsup7s73b9ntu0"
                }
              ]
            }
          ]
        },
        {
          "name": "Spectacles \u0026 Concerts",
          "slug": "spectacles-concerts",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genres",
              "slug": "genres",
              "tags": [
                {
                  "name": "Concerts",
                  "id": ""
                },
                {
                  "name": "Festivals",
                  "id": ""
                },
                {
                  "name": "Humour / One-man show",
                  "id": ""
                },
                {
                  "name": "Opéra / Ballet",
                  "id": ""
                },
                {
                  "name": "Théâtre",
                  "id": ""
                }
              ]
            },
            {
              "name": "Qualité / Résolution",
              "slug": "qualit-r-solution",
              "tags": [
                {
                  "name": "1080p (Full HD)",
                  "id": ""
                },
                {
                  "name": "2160p (4K)",
                  "id": ""
                },
                {
                  "name": "4320p (8K)",
                  "id": ""
                },
                {
                  "name": "720p (HD)",
                  "id": ""
                },
                {
                  "name": "SD",
                  "id": ""
                }
              ]
            },
            {
              "name": "Codec vidéo",
              "slug": "codec-vid-o",
              "tags": [
                {
                  "name": "AV1",
                  "id": ""
                },
                {
                  "name": "AVC/H264/x264",
                  "id": ""
                },
                {
                  "name": "HEVC/H265/x265",
                  "id": ""
                },
                {
                  "name": "MPEG",
                  "id": ""
                },
                {
                  "name": "VC-1",
                  "id": ""
                },
                {
                  "name": "VCC/H266/x266",
                  "id": ""
                },
                {
                  "name": "VP9",
                  "id": ""
                }
              ]
            },
            {
              "name": "Caractéristiques vidéo",
              "slug": "caract-ristiques-vid-o",
              "tags": [
                {
                  "name": "10 bits",
                  "id": ""
                },
                {
                  "name": "3D",
                  "id": ""
                },
                {
                  "name": "Dolby Vision",
                  "id": ""
                },
                {
                  "name": "HDR",
                  "id": ""
                },
                {
                  "name": "HDR10+",
                  "id": ""
                },
                {
                  "name": "HLG",
                  "id": ""
                },
                {
                  "name": "IMAX",
                  "id": ""
                },
                {
                  "name": "SDR",
                  "id": ""
                }
              ]
            },
            {
              "name": "Source / Type",
              "slug": "source-type",
              "tags": [
                {
                  "name": "4KLight",
                  "id": ""
                },
                {
                  "name": "BluRay",
                  "id": ""
                },
                {
                  "name": "DVDRip",
                  "id": ""
                },
                {
                  "name": "FULL Disc",
                  "id": ""
                },
                {
                  "name": "HDLight",
                  "id": ""
                },
                {
                  "name": "REMUX",
                  "id": ""
                },
                {
                  "name": "TV",
                  "id": ""
                },
                {
                  "name": "WEB-DL",
                  "id": ""
                },
                {
                  "name": "WEBRip",
                  "id": ""
                }
              ]
            },
            {
              "name": "Codec audio",
              "slug": "codec-audio",
              "tags": [
                {
                  "name": "AAC",
                  "id": ""
                },
                {
                  "name": "AC3",
                  "id": ""
                },
                {
                  "name": "AC4",
                  "id": ""
                },
                {
                  "name": "Autres",
                  "id": ""
                },
                {
                  "name": "DTS",
                  "id": ""
                },
                {
                  "name": "DTS-HD HR",
                  "id": ""
                },
                {
                  "name": "DTS-HD MA",
                  "id": ""
                },
                {
                  "name": "DTS:X",
                  "id": ""
                },
                {
                  "name": "E-AC3",
                  "id": ""
                },
                {
                  "name": "E-AC3 Atmos",
                  "id": ""
                },
                {
                  "name": "FLAC",
                  "id": ""
                },
                {
                  "name": "HE-AAC",
                  "id": ""
                },
                {
                  "name": "MP3",
                  "id": ""
                },
                {
                  "name": "Opus",
                  "id": ""
                },
                {
                  "name": "PCM",
                  "id": ""
                },
                {
                  "name": "TrueHD",
                  "id": ""
                },
                {
                  "name": "TrueHD Atmos",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues audio",
              "slug": "langues-audio",
              "tags": [
                {
                  "name": "Autre Langue",
                  "id": ""
                },
                {
                  "name": "Chinois",
                  "id": ""
                },
                {
                  "name": "English",
                  "id": ""
                },
                {
                  "name": "French",
                  "id": ""
                },
                {
                  "name": "Italian",
                  "id": ""
                },
                {
                  "name": "Japanese",
                  "id": ""
                },
                {
                  "name": "Korean",
                  "id": ""
                },
                {
                  "name": "MULTI",
                  "id": ""
                },
                {
                  "name": "Sans Dialogue",
                  "id": ""
                },
                {
                  "name": "Spanish",
                  "id": ""
                },
                {
                  "name": "VFF",
                  "id": ""
                },
                {
                  "name": "VFQ",
                  "id": ""
                }
              ]
            },
            {
              "name": "Sous-titres",
              "slug": "sous-titres",
              "tags": [
                {
                  "name": "Autres sous-titres",
                  "id": ""
                },
                {
                  "name": "ENG",
                  "id": ""
                },
                {
                  "name": "FR",
                  "id": ""
                },
                {
                  "name": "VFF Sous-Titres",
                  "id": ""
                },
                {
                  "name": "VFQ Sous-Titres",
                  "id": ""
                }
              ]
            },
            {
              "name": "Extension",
              "slug": "extension",
              "tags": [
                {
                  "name": "Autres Extensions",
                  "id": ""
                },
                {
                  "name": "AVI",
                  "id": ""
                },
                {
                  "name": "ISO",
                  "id": ""
                },
                {
                  "name": "MKV",
                  "id": ""
                },
                {
                  "name": "MP4",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Sports",
          "slug": "sports",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genres",
              "slug": "genres",
              "tags": [
                {
                  "name": "Basket-ball",
                  "id": ""
                },
                {
                  "name": "Cyclisme",
                  "id": ""
                },
                {
                  "name": "F1",
                  "id": ""
                },
                {
                  "name": "Foot",
                  "id": ""
                },
                {
                  "name": "Football américain",
                  "id": ""
                },
                {
                  "name": "Golf",
                  "id": ""
                },
                {
                  "name": "Handball",
                  "id": ""
                },
                {
                  "name": "Hockey sur glace",
                  "id": ""
                },
                {
                  "name": "Rugby",
                  "id": ""
                },
                {
                  "name": "Tennis",
                  "id": ""
                }
              ]
            },
            {
              "name": "Qualité / Résolution",
              "slug": "qualit-r-solution",
              "tags": [
                {
                  "name": "1080p (Full HD)",
                  "id": ""
                },
                {
                  "name": "2160p (4K)",
                  "id": ""
                },
                {
                  "name": "4320p (8K)",
                  "id": ""
                },
                {
                  "name": "720p (HD)",
                  "id": ""
                },
                {
                  "name": "SD",
                  "id": ""
                }
              ]
            },
            {
              "name": "Codec vidéo",
              "slug": "codec-vid-o",
              "tags": [
                {
                  "name": "AV1",
                  "id": ""
                },
                {
                  "name": "AVC/H264/x264",
                  "id": ""
                },
                {
                  "name": "HEVC/H265/x265",
                  "id": ""
                },
                {
                  "name": "MPEG",
                  "id": ""
                },
                {
                  "name": "VC-1",
                  "id": ""
                },
                {
                  "name": "VCC/H266/x266",
                  "id": ""
                },
                {
                  "name": "VP9",
                  "id": ""
                }
              ]
            },
            {
              "name": "Caractéristiques vidéo",
              "slug": "caract-ristiques-vid-o",
              "tags": [
                {
                  "name": "10 bits",
                  "id": ""
                },
                {
                  "name": "3D",
                  "id": ""
                },
                {
                  "name": "Dolby Vision",
                  "id": ""
                },
                {
                  "name": "HDR",
                  "id": ""
                },
                {
                  "name": "HDR10+",
                  "id": ""
                },
                {
                  "name": "HLG",
                  "id": ""
                },
                {
                  "name": "IMAX",
                  "id": ""
                },
                {
                  "name": "SDR",
                  "id": ""
                }
              ]
            },
            {
              "name": "Source / Type",
              "slug": "source-type",
              "tags": [
                {
                  "name": "4KLight",
                  "id": ""
                },
                {
                  "name": "BluRay",
                  "id": ""
                },
                {
                  "name": "DVDRip",
                  "id": ""
                },
                {
                  "name": "FULL Disc",
                  "id": ""
                },
                {
                  "name": "HDLight",
                  "id": ""
                },
                {
                  "name": "REMUX",
                  "id": ""
                },
                {
                  "name": "TV",
                  "id": ""
                },
                {
                  "name": "WEB-DL",
                  "id": ""
                },
                {
                  "name": "WEBRip",
                  "id": ""
                }
              ]
            },
            {
              "name": "Codec audio",
              "slug": "codec-audio",
              "tags": [
                {
                  "name": "AAC",
                  "id": ""
                },
                {
                  "name": "AC3",
                  "id": ""
                },
                {
                  "name": "AC4",
                  "id": ""
                },
                {
                  "name": "Autres",
                  "id": ""
                },
                {
                  "name": "DTS",
                  "id": ""
                },
                {
                  "name": "DTS-HD HR",
                  "id": ""
                },
                {
                  "name": "DTS-HD MA",
                  "id": ""
                },
                {
                  "name": "DTS:X",
                  "id": ""
                },
                {
                  "name": "E-AC3",
                  "id": ""
                },
                {
                  "name": "E-AC3 Atmos",
                  "id": ""
                },
                {
                  "name": "FLAC",
                  "id": ""
                },
                {
                  "name": "HE-AAC",
                  "id": ""
                },
                {
                  "name": "MP3",
                  "id": ""
                },
                {
                  "name": "Opus",
                  "id": ""
                },
                {
                  "name": "PCM",
                  "id": ""
                },
                {
                  "name": "TrueHD",
                  "id": ""
                },
                {
                  "name": "TrueHD Atmos",
                  "id": ""
                }
              ]
            },
            {
              "name": "Langues audio",
              "slug": "langues-audio",
              "tags": [
                {
                  "name": "Autre Langue",
                  "id": ""
                },
                {
                  "name": "Chinois",
                  "id": ""
                },
                {
                  "name": "English",
                  "id": ""
                },
                {
                  "name": "French",
                  "id": ""
                },
                {
                  "name": "Italian",
                  "id": ""
                },
                {
                  "name": "Japanese",
                  "id": ""
                },
                {
                  "name": "Korean",
                  "id": ""
                },
                {
                  "name": "MULTI",
                  "id": ""
                },
                {
                  "name": "Sans Dialogue",
                  "id": ""
                },
                {
                  "name": "Spanish",
                  "id": ""
                },
                {
                  "name": "VFF",
                  "id": ""
                },
                {
                  "name": "VFQ",
                  "id": ""
                }
              ]
            },
            {
              "name": "Sous-titres",
              "slug": "sous-titres",
              "tags": [
                {
                  "name": "Autres sous-titres",
                  "id": ""
                },
                {
                  "name": "ENG",
                  "id": ""
                },
                {
                  "name": "FR",
                  "id": ""
                },
                {
                  "name": "VFF Sous-Titres",
                  "id": ""
                },
                {
                  "name": "VFQ Sous-Titres",
                  "id": ""
                }
              ]
            },
            {
              "name": "Extension",
              "slug": "extension",
              "tags": [
                {
                  "name": "Autres Extensions",
                  "id": ""
                },
                {
                  "name": "AVI",
                  "id": ""
                },
                {
                  "name": "ISO",
                  "id": ""
                },
                {
                  "name": "MKV",
                  "id": ""
                },
                {
                  "name": "MP4",
                  "id": ""
                }
              ]
            }
          ]
        },
        {
          "name": "Séries TV",
          "slug": "series",
          "id": "cmjoyv2dg00067ery8m6c3q8h",
          "emplacementsouscategorie": null,
          "caracteristiques": [
            {
              "name": "Genres",
              "slug": "genres",
              "tags": [
                {
                  "name": "Action",
                  "id": "cmjudwh76000guyrus1jc9jxs"
                },
                {
                  "name": "Animation",
                  "id": "cmjudwhrd000iuyruhm4bd82d"
                },
                {
                  "name": "Aventure",
                  "id": "cmjudwhi0000huyru6x3xirry"
                },
                {
                  "name": "Biopic",
                  "id": "cmjudwl30000tuyrug7nl1dih"
                },
                {
                  "name": "Comédie",
                  "id": "cmjudwi0v000juyrujjoqy8ya"
                },
                {
                  "name": "Courts-métrages",
                  "id": "cmjudwld7000uuyrusigq41k6"
                },
                {
                  "name": "Documentaire",
                  "id": "cmjudwiay000kuyruuh6y6sux"
                },
                {
                  "name": "Drame",
                  "id": "cmjudwilv000luyruldwrewad"
                },
                {
                  "name": "Émission TV",
                  "id": "cmjudwteq001juyruql1f2o18"
                },
                {
                  "name": "Fantastique",
                  "id": "cmjudwiw7000muyruiwezf9u9"
                },
                {
                  "name": "Guerre",
                  "id": "cmjudwkto000suyruryh2tokp"
                },
                {
                  "name": "Historique",
                  "id": "cmjudwllt000vuyruryo0oj20"
                },
                {
                  "name": "Horreur",
                  "id": "cmjudwj5c000nuyruu32px04l"
                },
                {
                  "name": "Policier / Thriller",
                  "id": "cmjudwjf9000ouyrud0gadn2q"
                },
                {
                  "name": "Romance",
                  "id": "cmjudwkdw000ruyruwu1f7ps7"
                },
                {
                  "name": "Science-fiction",
                  "id": "cmjudwjri000puyru2t0vy9w9"
                },
                {
                  "name": "Sport",
                  "id": "cmjudwt42001iuyru9xod9658"
                },
                {
                  "name": "Western",
                  "id": "cmjudwk3e000quyrumxlv4nxg"
                }
              ]
            },
            {
              "name": "Qualité / Résolution",
              "slug": "qualit-r-solution",
              "tags": [
                {
                  "name": "1080p (Full HD)",
                  "id": "44d9f0e4-5ff2-4c6f-9e17-ec2e33acf448"
                },
                {
                  "name": "2160p (4K)",
                  "id": "e42ac76f-d1b0-4ea4-a3da-f332af0f8f4c"
                },
                {
                  "name": "4320p (8K)",
                  "id": "20d145d3-db89-4996-933c-48cf51ee5b6b"
                },
                {
                  "name": "720p (HD)",
                  "id": "86e547c3-7416-46c5-a759-fc1a5d32cc23"
                },
                {
                  "name": "SD",
                  "id": "cmjudwm65000xuyrubkb2ztkf"
                }
              ]
            },
            {
              "name": "Codec vidéo",
              "slug": "codec-vid-o",
              "tags": [
                {
                  "name": "AV1",
                  "id": "cmjudwnvj0010uyrusdqdsydj"
                },
                {
                  "name": "AVC/H264/x264",
                  "id": "cmjoyv2id000u7eryugoe1bee"
                },
                {
                  "name": "HEVC/H265/x265",
                  "id": "cmjoyv2ig000v7eryc9hf1hsa"
                },
                {
                  "name": "MPEG",
                  "id": "5c9bb557-4fd7-488a-8ce7-83c26c7049f2"
                },
                {
                  "name": "VC-1",
                  "id": "168bff0e-8649-4fb4-86d9-87006c2aa511"
                },
                {
                  "name": "VCC/H266/x266",
                  "id": "f77bbe91-82c6-440d-a302-c3015edc19a8"
                },
                {
                  "name": "VP9",
                  "id": "34d10b63-fcd7-4b12-91df-48cde5cb63c0"
                }
              ]
            },
            {
              "name": "Caractéristiques vidéo",
              "slug": "caract-ristiques-vid-o",
              "tags": [
                {
                  "name": "10 bits",
                  "id": "d3827daf-6d53-4d61-8e88-5a5b33ddd85f"
                },
                {
                  "name": "3D",
                  "id": "d5eg09psup7s739bto70"
                },
                {
                  "name": "Dolby Vision",
                  "id": "862f0aaf-e08d-492f-bf2d-9806c74b676a"
                },
                {
                  "name": "HDR",
                  "id": "4e2f5500-05f9-4f35-b273-233abc8ff991"
                },
                {
                  "name": "HDR10+",
                  "id": "da699e63-de34-4c75-8d65-243d8eb51151"
                },
                {
                  "name": "HLG",
                  "id": "79771165-c073-490d-81e3-a408b01cfaa6"
                },
                {
                  "name": "IMAX",
                  "id": "d5gh4ohsup7s73b5irt0"
                },
                {
                  "name": "SDR",
                  "id": "2a3d5386-ba5f-4cce-b957-1456a6a938da"
                }
              ]
            },
            {
              "name": "Source / Type",
              "slug": "source-type",
              "tags": [
                {
                  "name": "4KLight",
                  "id": "d77a8e79-0035-4df8-8cef-8311a1ea1919"
                },
                {
                  "name": "BluRay",
                  "id": "a063935d-adc5-4e69-af04-54f318a80771"
                },
                {
                  "name": "DVDRip",
                  "id": "48c926bc-8fb6-4163-99b7-c745aba4f1ed"
                },
                {
                  "name": "FULL Disc",
                  "id": "fca2b774-3587-426c-8b7a-08e0fe51f2c6"
                },
                {
                  "name": "HDLight",
                  "id": "a1c90ed0-a4b9-4444-bd8d-a4fce8dc5caf"
                },
                {
                  "name": "REMUX",
                  "id": "f4e1b729-bc12-4957-94ee-53d10bfbf63a"
                },
                {
                  "name": "TV",
                  "id": "f8d828b1-1d48-443f-8749-dd4af5cba373"
                },
                {
                  "name": "WEB-DL",
                  "id": "7bd8b291-6e18-4322-9c35-b3470c90039e"
                },
                {
                  "name": "WEBRip",
                  "id": "86413ec8-af63-4fe2-8c88-cb56ec1b176c"
                }
              ]
            },
            {
              "name": "Codec audio",
              "slug": "codec-audio",
              "tags": [
                {
                  "name": "AAC",
                  "id": "cmjudwo5h0011uyruofwyb2cn"
                },
                {
                  "name": "AC3",
                  "id": "cmjudwodw0012uyrut8jh456x"
                },
                {
                  "name": "AC4",
                  "id": "d5e5hc1sup7s73ag6eig"
                },
                {
                  "name": "Autres",
                  "id": "d5e5hnhsup7s73ecfpl0"
                },
                {
                  "name": "DTS",
                  "id": "cmjudwomi0013uyruthukb2pf"
                },
                {
                  "name": "DTS-HD HR",
                  "id": "5e8475fe-db59-4dd8-84fe-7186eaba134d"
                },
                {
                  "name": "DTS-HD MA",
                  "id": "f87c6b8e-6edf-4dbd-b642-205d1060c0d1"
                },
                {
                  "name": "DTS:X",
                  "id": "ec308312-a650-4b00-b011-e624c8779939"
                },
                {
                  "name": "E-AC3",
                  "id": "3ab66a11-d74c-49b7-a3f9-2efc0edfefb2"
                },
                {
                  "name": "E-AC3 Atmos",
                  "id": "38ea7ff5-bd03-4b41-af5c-211792c0d1e8"
                },
                {
                  "name": "FLAC",
                  "id": "d5e6q01sup7s738q0tsg"
                },
                {
                  "name": "HE-AAC",
                  "id": "c1f82902-ad95-4a1c-abdb-3dd8f12b1f40"
                },
                {
                  "name": "MP3",
                  "id": "d5e6q2psup7s738q0tt0"
                },
                {
                  "name": "Opus",
                  "id": "ebad8b2b-0f23-4a3a-a741-9098e61d5f46"
                },
                {
                  "name": "PCM",
                  "id": "4b6f7605-3aa3-411e-8bb4-fd6fcbc9c921"
                },
                {
                  "name": "TrueHD",
                  "id": "c45a5dbb-aad2-461c-96d2-ee4cee3f5d5b"
                },
                {
                  "name": "TrueHD Atmos",
                  "id": "963a0227-d20b-4878-aa84-1deac1de4479"
                }
              ]
            },
            {
              "name": "Langues audio",
              "slug": "langues-audio",
              "tags": [
                {
                  "name": "Autre Langue",
                  "id": "d5elgi9sup7s73emca6g"
                },
                {
                  "name": "Chinois",
                  "id": "d5eldm1sup7s7387o7ng"
                },
                {
                  "name": "English",
                  "id": "d5e5kn1sup7s73b1q4pg"
                },
                {
                  "name": "French",
                  "id": "2d45b5d1-2dfe-4de5-b0fe-e4a08132d06a"
                },
                {
                  "name": "Italian",
                  "id": "d5e5jh1sup7s73ecfppg"
                },
                {
                  "name": "Japanese",
                  "id": "d5e5jchsup7s73ecfpog"
                },
                {
                  "name": "Korean",
                  "id": "d5e5jepsup7s73f6hn9g"
                },
                {
                  "name": "MULTI",
                  "id": "cmjoyv2i0000q7eryz0dnb7f9"
                },
                {
                  "name": "Sans Dialogue",
                  "id": "d5hb34pua6hc739gjva0"
                },
                {
                  "name": "Spanish",
                  "id": "d5e5mhhsup7s73cq0sr0"
                },
                {
                  "name": "VFF",
                  "id": "5cc1e14d-b23a-42c5-956e-3ad9e1529b3c"
                },
                {
                  "name": "VFQ",
                  "id": "b4661137-ed7a-4742-8ebc-40fcb6fa8f8d"
                }
              ]
            },
            {
              "name": "Sous-titres",
              "slug": "sous-titres",
              "tags": [
                {
                  "name": "Autres sous-titres",
                  "id": "d5elg3psup7s73emvf3g"
                },
                {
                  "name": "ENG",
                  "id": "d5e6lu9sup7s73a9mimg"
                },
                {
                  "name": "FR",
                  "id": "cmjudwg9q000euyruscojd2q4"
                },
                {
                  "name": "VFF Sous-Titres",
                  "id": "d5elf41sup7s73emc85g"
                },
                {
                  "name": "VFQ Sous-Titres",
                  "id": "d5elilhsup7s73avoesg"
                }
              ]
            },
            {
              "name": "Extension",
              "slug": "extension",
              "tags": [
                {
                  "name": "Autres Extensions",
                  "id": "d5fuf51sup7s73bbcmq0"
                },
                {
                  "name": "AVI",
                  "id": "d5fueopsup7s73cg5qo0"
                },
                {
                  "name": "ISO",
                  "id": "d5fuevpsup7s739du1l0"
                },
                {
                  "name": "MKV",
                  "id": "d5fuer9sup7s73eq24s0"
                },
                {
                  "name": "MP4",
                  "id": "d5fuelpsup7s73b9ntu0"
                }
              ]
            }
          ]
        }
      ],
      "caracteristiques": null
    }
  ]
}`
