<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

# Agent Guidelines for HelloCommit

## Project Overview

This is a Next.js 16 (App Router) project with React 19, TypeScript, and Tailwind CSS v4. It uses Base UI components styled with shadcn/ui conventions, Lucide icons, and class-variance-authority for component variants.

## Commands

```bash
# Development
npm run dev           # Start dev server at http://localhost:3000

# Production
npm run build         # Build for production
npm run start         # Start production server

# Code Quality
npm run lint          # Run ESLint on all files
```

## Code Style Guidelines

### TypeScript

- Use strict TypeScript (`"strict": true` in tsconfig.json)
- Prefer explicit types for function parameters and return values
- Use `type` for type aliases, `interface` for object shapes
- Avoid `any` — use `unknown` when type is truly unknown
- Use optional chaining (`?.`) and nullish coalescing (`??`) for null/undefined handling
- Import types directly: `import { type ClassValue } from "clsx"`

### Imports

- Use absolute path imports with `@/` prefix (configured in tsconfig.json)
- Order imports: external packages → internal modules → relative imports → types
- Group imports without blank lines; use separate blocks for types
- Example:
  ```typescript
  import { clsx, type ClassValue } from "clsx"
  import { twMerge } from "tailwind-merge"
  
  import { Button } from "@/components/ui/button"
  import { cn } from "@/lib/utils"
  ```

### Naming Conventions

- **Files**: kebab-case for config files, PascalCase for React components, camelCase for utilities
  - Components: `button.tsx`, `user-profile.tsx`
  - Utils: `utils.ts`, `format-date.ts`
- **Functions**: camelCase, use verb prefixes for event handlers (`handleClick`)
- **Components**: PascalCase, descriptive names
- **CSS classes**: Use Tailwind utility classes; avoid custom CSS unless necessary

### React Patterns

- Use Server Components by default (no "use client" directive)
- Add `"use client"` only when using client-side features (hooks, event handlers, browser APIs)
- Use `Readonly<{...}>` type for props in server components
- Prefer functional components with explicit return types for clarity
- Export components both as default and named exports when useful (e.g., `export { Button, buttonVariants }`)

### Tailwind CSS v4

- Use `@import "tailwindcss"` syntax (v4)
- CSS variables go in `@theme inline {}` block
- Use oklch color format for design tokens
- Use Tailwind's `@apply` sparingly
- Dark mode via `.dark` class on root element

### Component Structure (shadcn/ui style)

```typescript
"use client"  // if needed

import { Component } from "external-package"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

// CVA variants first
const componentVariants = cva("base-classes", {
  variants: { ... },
  defaultVariants: { ... },
})

// Component function
function Component({
  className,
  variant = "default",
  ...props
}: ComponentProps & VariantProps<typeof componentVariants>) {
  return (
    <PrimitiveComponent
      data-slot="component-name"
      className={cn(componentVariants({ variant, className }))}
      {...props}
    />
  )
}

export { Component, componentVariants }
```

### Error Handling

- Let TypeScript errors surface during build
- Use `aria-invalid` for form validation states
- Use semantic HTML for error messages
- No try/catch unless handling specific recoverable errors

### File Organization

```
/app              # Next.js App Router pages and layouts
  /page.tsx       # Home page
  /layout.tsx     # Root layout with metadata
  /globals.css    # Tailwind imports and CSS variables
/components
  /ui/            # shadcn-style base components
/lib
  /utils.ts       # Shared utilities (cn function)
/public           # Static assets
```

### Path Aliases

- `@/*` → project root
- `@/components/*` → components directory
- `@/lib/*` → lib directory
- `@/components/ui` → UI component library

## Important Notes

- This is Next.js 16 with App Router — APIs may differ from earlier versions
- Use Base UI components (`@base-ui/react`) as primitives
- No CSS modules — use Tailwind exclusively for styling
- No dark mode CSS variables beyond what's in globals.css
- Test changes with `npm run lint` before committing
