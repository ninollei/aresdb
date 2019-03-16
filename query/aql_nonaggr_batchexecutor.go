package query

import (
	"github.com/uber/aresdb/memutils"
	queryCom "github.com/uber/aresdb/query/common"
	"unsafe"
)

// NonAggrBatchExecutorImpl is batch executor implementation for non-aggregation query
type NonAggrBatchExecutorImpl struct {
	*BatchExecutorImpl
}

// reduce function for NonAggrBatchExecutorImpl will expand the result instead from baseCount
func (e *NonAggrBatchExecutorImpl) reduce() {
  // call expand here
}

// project for non-aggregation query will only calculate the selected columns
// measure calculation, reduce will be skipped, once the generated result reaches limit, it will return and cancel all other ongoing processing.
func (e *NonAggrBatchExecutorImpl) project() {
	// Prepare for dimension and measure evaluation.
	e.prepareForDimEval(e.qc.OOPK.DimRowBytes, e.qc.OOPK.NumDimsPerDimWidth, e.stream)

	e.qc.reportTimingForCurrentBatch(e.stream, &e.start, prepareForDimAndMeasureTiming)

	e.evalDimensions()

	// wait for stream to clean up non used buffer before final aggregation
	memutils.WaitForCudaStream(e.stream, e.qc.Device)
	e.qc.OOPK.currentBatch.cleanupBeforeAggregation()
}

func (e *NonAggrBatchExecutorImpl) prepareForDimEval(
	dimRowBytes int, numDimsPerDimWidth queryCom.DimCountsPerDimWidth, stream unsafe.Pointer) {

	bc := &e.qc.OOPK.currentBatch
	if bc.resultCapacity == 0 {
		bc.resultCapacity = e.qc.Query.Limit
		bc.dimensionVectorD = [2]devicePointer{
			deviceAllocate(bc.resultCapacity*dimRowBytes, bc.device),
			deviceAllocate(bc.resultCapacity*dimRowBytes, bc.device),
		}
	}
	if bc.size + bc.resultSize > bc.resultCapacity {
		bc.size = bc.resultCapacity - bc.resultSize
	}
}


