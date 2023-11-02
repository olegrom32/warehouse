package packager

import (
	"fmt"
	"math"
	"sort"
)

type Service struct {
	sizes []int
}

func NewService(sizes []int) *Service {
	used := map[int]struct{}{}
	preparedSizes := []int{}

	// Remove duplicates and sort the sizes so the algo works properly
	for _, size := range sizes {
		if _, ok := used[size]; ok {
			continue
		}

		preparedSizes = append(preparedSizes, size)
		used[size] = struct{}{}
	}

	sort.Slice(preparedSizes, func(i, j int) bool {
		return preparedSizes[i] > preparedSizes[j]
	})

	return &Service{sizes: preparedSizes}
}

func (s *Service) Package(totalItems int) map[int]int {
	// Some base cases
	if totalItems == 0 || len(s.sizes) == 0 {
		return nil
	}

	// Use the cache to store possible best options, so we can then return the best one after we've
	// calculated all the options
	cache := map[string]map[int]int{}

	// Two pointers to keep track of if a given option is better than previous
	bestItems := math.MaxInt
	bestBoxes := math.MaxInt

	// Keep track of how many items are fulfilled using bigger boxes, so we can try to fulfil the remained
	// using smaller boxes
	accumulatedItems := 0
	accumulatedBoxes := map[int]int{}

	// Iterate through all the box sizes once
	for _, size := range s.sizes {
		// Calculate the amount of boxes that fit into the requested item amount
		numOfBoxes := (totalItems - accumulatedItems) / size

		// Save the possible best option to the cache:
		// 1. if the total amount of items, that we will ship with these boxes is less (better)
		// 		than what we have already calculated
		if accumulatedItems+(size*(numOfBoxes+1)) < bestItems {
			bestItems = accumulatedItems + (size * (numOfBoxes + 1))
			bestBoxes = math.MaxInt
		}

		// 2. if the number of boxes needed to fulfil the request is less (better)
		//		than what we have already calculated
		if len(accumulatedBoxes)+numOfBoxes+1 < bestBoxes {
			bestBoxes = len(accumulatedBoxes) + numOfBoxes + 1

			// "totalItems:numOfBoxes"
			cacheKey := fmt.Sprintf("%d:%d", bestItems, bestBoxes)
			// If all the conditions match - the current option is currently the best possible option, so
			// we should cache it until we process all the rest of options
			cache[cacheKey] = map[int]int{}
			for k, v := range accumulatedBoxes {
				cache[cacheKey][k] = v
			}
			cache[cacheKey][size] = numOfBoxes + 1
		}

		if numOfBoxes > 0 {
			// Add the whole number of boxes of this size to the total and try to fulfil the remainder of items
			// using smaller boxes
			accumulatedItems += size * numOfBoxes
			accumulatedBoxes[size] = numOfBoxes
		}
	}

	// A special case when we processed all the options and the very last one turned out to be the best
	if accumulatedItems == totalItems {
		return accumulatedBoxes
	}

	// Return the best option
	return cache[fmt.Sprintf("%d:%d", bestItems, bestBoxes)]
}
