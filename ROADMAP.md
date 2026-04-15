# HelloCommit — Roadmap

Roadmap pragmatique pour HelloCommit, ordonnée par priorité : livrer vite, scaler ensuite.

**Stack cible :** Next.js (front) + Go (API) + SQLite (DB) + Litestream (backup) sur un VPS Hetzner, avec Vercel pour le front.

---

## Phase 1 — MVP solide (semaines 1-2)

**Objectif : que ça marche bien pour 10 users beta.**

- [ ] Schéma DB propre : tables `users`, `repos`, `user_stars`, `issues` (avec `last_scanned_at`, `is_open`, `activity_score`)
- [ ] OAuth GitHub côté Next.js, stockage du token chiffré côté Go
- [ ] Endpoint Go `POST /sync-stars` : récupère les starred repos de l'user et les insère dans `user_stars` + `repos`
- [ ] Endpoint Go `GET /issues` : retourne les GFI des repos starrés par l'user
- [ ] Worker de scan basique en goroutine : boucle qui pick les repos `last_scanned_at < now - 1h`
- [ ] Page Next.js qui liste les issues avec lien vers GitHub

---

## Phase 2 — Robustesse (semaines 3-4)

**Objectif : ne pas se faire bannir par GitHub et ne pas perdre de données.**

- [ ] **ETags / conditional requests** sur les appels GitHub (énorme gain rate limit)
- [ ] Gestion propre du rate limit : lire les headers `X-RateLimit-Remaining`, backoff exponentiel
- [ ] **Litestream** vers Cloudflare R2 pour les backups continus
- [ ] Logs structurés (`slog` en Go) + endpoint `/health`
- [ ] Couche d'abstraction DB (interface Go `Repository`) pour préparer une éventuelle migration Postgres
- [ ] Gestion des labels variés : `good first issue`, `good-first-issue`, `beginner`, configurable

---

## Phase 3 — Qualité produit (semaines 5-6)

**Objectif : que les issues affichées soient vraiment utiles.**

- [ ] **Filtrage des issues mortes** : exclure les issues sans activité depuis X mois, ou déjà assignées
- [ ] Score de pertinence : récence, langage, taille du repo, activité récente
- [ ] Filtres UI : par langage, par techno, par taille de projet
- [ ] Détection des issues "claimed" via les commentaires (« I'll work on this »)
- [ ] Refresh manuel pour l'user (avec rate limit côté toi)

---

## Phase 4 — Croissance (quand tu approches 1k users)

**Objectif : préparer le scaling sans casser la prod.**

- [ ] Cache in-memory (ristretto) pour les requêtes fréquentes
- [ ] Métriques basiques : Prometheus + Grafana Cloud free tier, ou juste des logs
- [ ] Worker de scan séparé du serveur API (toujours même VPS, mais process distinct)
- [ ] Notifications email/web push pour nouvelles GFI matchant les centres d'intérêt
- [ ] Page publique « trending GFIs » (sans auth) — bon pour le SEO

---

## Phase 5 — Scale (si tu dépasses 1k users actifs)

**Objectif : passer le mur du single-VPS.**

- [ ] Migration SQLite → Postgres (Neon)
- [ ] Worker de scan sur sa propre machine
- [ ] API Go stateless, déployable en plusieurs instances
- [ ] CDN devant Next.js (déjà gratuit avec Vercel)
- [ ] Queue dédiée pour les scans (NATS ou River en Go, pas besoin de Kafka)

---

## En parallèle, dès le début

- Un fichier `ARCHITECTURE.md` dans le repo qui explique les choix → utile pour toi dans 6 mois
- Tests d'intégration sur les endpoints critiques (`sync-stars`, scan worker)
- Un script de seed pour dev local avec des fake data