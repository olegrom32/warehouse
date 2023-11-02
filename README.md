# warehouse

## Structure

- main.go - the entrypoint + REST API
- internal/packager - an implementation of the packager

Please see code comments for more details.

## Run the app

Simply run to run with default box sizes
```shell
$ docker-compose up
```

To configure box sizes (the format for the env variable is a json array of integers):
```shell
SIZES="[5000,2000,100]" docker-compose up
```

Services exposed:

> Super simple UI: http://localhost

> REST API: http://localhost:8080/order

> Docs: http://localhost:8088

## Algorithm

Based on the greedy algorithm.

The algorithm is running through the list of available box sizes and tries to find optimal solutions.
The solution with least total items and least number of boxes wins.

Breakdown of some cases (for box sizes `[5000, 2000, 1000, 500, 250]`):

**Target: 751**

1. try box 5000 (5000 items total, 1 box needed, possible winner)
2. try box 2000 (2000 items total - better, 1 box needed, possible winner)
3. try box 1000 (1000 items total - better, 1 box needed, possible winner)
4. try boxes 500+500 (1000 items total, 2 boxes needed - worse, not winner)
5. try boxes 500+250+250 (1000 items total, 3 boxes needed - worse, not winner)
6. we ran out of box sizes, return the last best (from step 3)

**Target: 12001**

1. try box 5000 (15000 items total, 3 boxes needed, possible winner)
2. try boxes 5000+5000 and the remainder try with boxes 2000+2000 (14000 items total - better, 4 boxes needed, possible winner)
3. try boxes 5000+5000+2000 and the remainder try with box 1000 (13000 items total - better, 4 boxes needed, possible winner)
4. try boxes 5000+5000+2000 and the remainder try with box 500 (12500 items total - better,  4 boxes needed, possible winner)
5. try boxes 5000+5000+2000 and the remainder try box 250 (12250 items total - better, 4 boxes total, possible winner)
6. we ran out of box sizes, return the last best (from step 5)

## Tests

`internal/packager` has 100% coverage.

Run all tests:
```shell
$ make test
```
