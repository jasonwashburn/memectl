## Why

Managed memes can be listed but not inspected individually. Users need a
kubectl-shaped way to review the complete locally stored record, including its
ordered caption text and hosted URLs, without contacting Imgflip.

## What Changes

- Add `memectl describe meme <name>` to display one locally managed meme in a
  labeled detail view.
- Add `memectl desc meme <name>` as the conventional short alias.
- Display the stored name, template ID, ordered caption texts, image URL, page
  URL, and exact UTC creation timestamp.
- Validate the requested name and return actionable errors for invalid or
  missing local records.
- Document that describe reads only the local inventory and does not contact
  Imgflip.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `meme-inventory`: Add detailed inspection of a named locally managed meme.

## Impact

- Affects Cobra command registration and a new describe command implementation.
- Reuses the existing local inventory and adds no Imgflip API calls,
  credentials, persisted fields, or dependencies.
- Updates command tests and README usage documentation.
