# ── Backend (Go / Fiber) ──────────────────────────────────────────────────────
docker_build(
  'hellocommit-api',
  context='.',
  dockerfile='Dockerfile',
  only=[
    './cmd/',
    './internal/',
    './pkg/',
    './go.mod',
    './go.sum',
  ],
)

# ── Frontend (Next.js) ────────────────────────────────────────────────────────
docker_build(
  'hellocommit-frontend',
  context='./ui',
  dockerfile='./ui/Dockerfile',
  only=[
    './app/',
    './components/',
    './hooks/',
    './lib/',
    './types/',
    './public/',
    './auth.ts',
    './proxy.ts',
    './next.config.ts',
    './postcss.config.mjs',
    './tsconfig.json',
    'package.json',
    'package-lock.json',
  ],
)

# ── Helm chart ────────────────────────────────────────────────────────────────
values_files = ['./chart/values.yaml']
if os.path.exists('./chart/values.local.yaml'):
  values_files.append('./chart/values.local.yaml')

k8s_yaml(helm(
  './chart',
  name='hellocommit',
  values=values_files,
))

# ── Resources ─────────────────────────────────────────────────────────────────
k8s_resource(
  'hellocommit-api',
  port_forwards=['8080:8080'],
  labels=['backend'],
)

k8s_resource(
  'hellocommit-frontend',
  port_forwards=['3000:3000'],
  labels=['frontend'],
)
