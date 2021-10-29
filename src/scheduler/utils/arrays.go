/*
 * P2PFaaS - A framework for FaaS Load Balancing
 * Copyright (c) 2019. Gabriele Proietti Mattia <pm.gabriele@outlook.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package utils

func MaxOfArrayFloat(array []float64) (float64, int) {
	maxValue := array[0]
	maxIndex := 0
	for i, v := range array {
		if v > maxValue {
			maxIndex = i
			maxValue = v
		}
	}
	return maxValue, maxIndex
}

func MaxOfArrayUint(array []uint) (uint, int) {
	maxValue := array[0]
	maxIndex := 0
	for i, v := range array {
		if v > maxValue {
			maxIndex = i
			maxValue = v
		}
	}
	return maxValue, maxIndex
}

func MinOfArrayUint(array []uint) (uint, int) {
	minValue := array[0]
	minIndex := 0
	for i, v := range array {
		if v < minValue {
			minIndex = i
			minValue = v
		}
	}
	return minValue, minIndex
}

func SlotsAboveSpecificFreeSlots(slots []uint, threshold uint) []uint {
	var slotsBelow []uint
	for i, v := range slots {
		if v > threshold {
			slotsBelow = append(slotsBelow, uint(i))
		}
	}
	return slotsBelow
}

func LoadsBelowSpecificLoad(slots []uint, threshold uint) []uint {
	var loadsBelow []uint
	for i, v := range slots {
		if v <= threshold {
			loadsBelow = append(loadsBelow, uint(i))
		}
	}
	return loadsBelow
}
