---
description: Push changes to GitHub
---

1. Stage all changes
   // turbo
   ```bash
   git add .
   ```

2. Commit with descriptive message
   ```bash
   git commit -m "${COMMIT_MESSAGE:-update}"
   ```

3. Push to origin main
   // turbo
   ```bash
   git push origin main
   ```
