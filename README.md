# Match

Match can find the best partners for you customer in the housing market.

At the moment the app just tires to find the best match by the partner's average rating and distance to the customer (given a latitude and longitude for both the customer and the partner).
The distance between partners and customers is calculated using the [Haversine Formula](https://en.wikipedia.org/wiki/Haversine_formula). 

The app also verifies if the partner has the requested materials, and it completely ignores the category of the partner and any other data sent by the customer like square foot of material, etc.

For simplicity, each partner has categories and materials and these entities are not connected to each other.
Also, it wasn't given any friendly ID or code to any category or material.

## Run

To run the application run the following command (from the root directory):

```shell
make docker-up
```

To stop the application run the following command (from the root directory):

```shell
make docker-down
```

## Test

### Mocks

To run the unit tests you first need to generate the mocks. To do so, run the following command (from the root directory):

```shell
make mocks
```

### Unit tests

To run the unit tests run the following command (from the root directory):

```shell
make test_unit
```

### Integration tests

To run the integration tests run the following commands (from the root directory):

```shell
make docker-up
```

and then:

```shell
make test_integration
```

Note that changing the ports in the file `docker-compose.yml` or changing the scripts in `01-init.sql` (in the `/scripts/db/` directory) can make the integration tests fail.

## Lint

To lint the code run the following command (from the root directory):

```shell
make vet
```
