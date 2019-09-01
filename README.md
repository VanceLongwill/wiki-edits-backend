## Hatnote historical

This project stores a live feed of wikipedia edits from hatnote's api in a postgres db. This data is used to provide a net change (bytes) result for a given time period.

### Endpoints

- GET `/edits?langCode=en&from=1567272675&to=1567272682`

  ```json
  {
    "netChange": 12897
  }
  ```

  - **Params**:
  - `langCode`: "en" | "de"
  - `from`: unix timestamp
  - `to`: unix timestamp

## Running the project

- `docker-compose up --build`
