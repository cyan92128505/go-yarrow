package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-drift/drift/pkg/platform"
)

// DivinationRecord 單筆占卜記錄
type DivinationRecord struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	CreatedAt time.Time `json:"created_at"`
	Seed      int64     `json:"seed"`
	Original  []int     `json:"original"`  // Line values: 6,7,8,9
	Changed   []int     `json:"changed"`   // Line values after transformation
	MovingPos []int     `json:"moving_pos"` // 1-based positions of moving lines
	Interpret string    `json:"interpret"`  // 朱熹法解卦結果
}

// ToHexagram reconstructs the Hexagram from stored data.
func (r *DivinationRecord) ToHexagram() *Hexagram {
	hex := &Hexagram{
		Original:    make([]Line, len(r.Original)),
		Changed:     make([]Line, len(r.Changed)),
		MovingLines: r.MovingPos,
	}
	for i, v := range r.Original {
		hex.Original[i] = Line(v)
	}
	for i, v := range r.Changed {
		hex.Changed[i] = Line(v)
	}
	return hex
}

// NewRecordFromHexagram creates a record from a hexagram result.
func NewRecordFromHexagram(question string, seed int64, hex *Hexagram) *DivinationRecord {
	now := time.Now()
	original := make([]int, len(hex.Original))
	changed := make([]int, len(hex.Changed))
	for i, l := range hex.Original {
		original[i] = int(l)
	}
	for i, l := range hex.Changed {
		changed[i] = int(l)
	}

	return &DivinationRecord{
		ID:        fmt.Sprintf("%d", now.UnixNano()),
		Question:  question,
		CreatedAt: now,
		Seed:      seed,
		Original:  original,
		Changed:   changed,
		MovingPos: hex.MovingLines,
		Interpret: hex.InterpretZhuXi(),
	}
}

// RecordStore manages persistence of divination records.
type RecordStore struct {
	mu      sync.Mutex
	records []*DivinationRecord
	dataDir string
}

const recordsFileName = "divination_records.json"

// NewRecordStore creates a store and loads existing data.
func NewRecordStore() *RecordStore {
	s := &RecordStore{}
	s.load()
	return s
}

func (s *RecordStore) filePath() string {
	if s.dataDir != "" {
		return s.dataDir + "/" + recordsFileName
	}
	// Fallback: try to get app documents directory
	docsPath, err := platform.Storage.GetAppDirectory(platform.AppDirectoryDocuments)
	if err == nil && docsPath != "" {
		s.dataDir = docsPath
		return docsPath + "/" + recordsFileName
	}
	// Last resort
	return recordsFileName
}

func (s *RecordStore) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := platform.Storage.ReadFile(s.filePath())
	if err != nil {
		s.records = make([]*DivinationRecord, 0)
		return
	}

	var records []*DivinationRecord
	if err := json.Unmarshal(data, &records); err != nil {
		s.records = make([]*DivinationRecord, 0)
		return
	}
	s.records = records
}

func (s *RecordStore) save() error {
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}
	return platform.Storage.WriteFile(s.filePath(), data)
}

// Add inserts a record and persists.
func (s *RecordStore) Add(record *DivinationRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Prepend (newest first)
	s.records = append([]*DivinationRecord{record}, s.records...)
	return s.save()
}

// Delete removes a record by ID and persists.
func (s *RecordStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.records {
		if r.ID == id {
			s.records = append(s.records[:i], s.records[i+1:]...)
			return s.save()
		}
	}
	return nil
}

// All returns all records sorted by creation time (newest first).
func (s *RecordStore) All() []*DivinationRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]*DivinationRecord, len(s.records))
	copy(result, s.records)
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result
}

// Count returns the number of records.
func (s *RecordStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.records)
}

// FindByID returns a record by ID.
func (s *RecordStore) FindByID(id string) *DivinationRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.records {
		if r.ID == id {
			return r
		}
	}
	return nil
}
