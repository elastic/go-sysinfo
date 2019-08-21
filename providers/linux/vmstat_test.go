// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package linux

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawInput = `
nr_free_pages 50545
nr_zone_inactive_anon 66
nr_zone_active_anon 26799
nr_zone_inactive_file 31849
nr_zone_active_file 94164
nr_zone_unevictable 0
nr_zone_write_pending 7
nr_mlock 0
nr_page_table_pages 1225
nr_kernel_stack 2496
nr_bounce 0
nr_zspages 0
nr_free_cma 0
numa_hit 44470329
numa_miss 0
numa_foreign 0
numa_interleave 16296
numa_local 44470329
numa_other 0
nr_inactive_anon 66
nr_active_anon 26799
nr_inactive_file 31849
nr_active_file 94164
nr_unevictable 0
nr_slab_reclaimable 31763
nr_slab_unreclaimable 10329
nr_isolated_anon 0
nr_isolated_file 0
workingset_refault 302914
workingset_activate 108959
workingset_nodereclaim 6422
nr_anon_pages 26218
nr_mapped 8641
nr_file_pages 126182
nr_dirty 7
nr_writeback 0
nr_writeback_temp 0
nr_shmem 169
nr_shmem_hugepages 0
nr_shmem_pmdmapped 0
nr_anon_transparent_hugepages 0
nr_unstable 0
nr_vmscan_write 35
nr_vmscan_immediate_reclaim 9832
nr_dirtied 7188920
nr_written 6479005
nr_dirty_threshold 31736
nr_dirty_background_threshold 15848
pgpgin 17010697
pgpgout 27734292
pswpin 0
pswpout 0
pgalloc_dma 241378
pgalloc_dma32 45788683
pgalloc_normal 0
pgalloc_movable 0
allocstall_dma 0
allocstall_dma32 0
allocstall_normal 5
allocstall_movable 8
pgskip_dma 0
pgskip_dma32 0
pgskip_normal 0
pgskip_movable 0
pgfree 46085578
pgactivate 2475069
pgdeactivate 636658
pglazyfree 9426
pgfault 46777498
pgmajfault 19204
pglazyfreed 0
pgrefill 707817
pgsteal_kswapd 3798890
pgsteal_direct 1466
pgscan_kswapd 3868525
pgscan_direct 1483
pgscan_direct_throttle 0
zone_reclaim_failed 0
pginodesteal 1710
slabs_scanned 8348560
kswapd_inodesteal 3142001
kswapd_low_wmark_hit_quickly 541
kswapd_high_wmark_hit_quickly 332
pageoutrun 1492
pgrotated 29725
drop_pagecache 0
drop_slab 0
oom_kill 0
numa_pte_updates 0
numa_huge_pte_updates 0
numa_hint_faults 0
numa_hint_faults_local 0
numa_pages_migrated 0
pgmigrate_success 4539
pgmigrate_fail 156
compact_migrate_scanned 9331
compact_free_scanned 136266
compact_isolated 9407
compact_stall 2
compact_fail 0
compact_success 2
compact_daemon_wake 21
compact_daemon_migrate_scanned 8311
compact_daemon_free_scanned 107086
htlb_buddy_alloc_success 0
htlb_buddy_alloc_fail 0
unevictable_pgs_culled 19
unevictable_pgs_scanned 0
unevictable_pgs_rescued 304
unevictable_pgs_mlocked 304
unevictable_pgs_munlocked 304
unevictable_pgs_cleared 0
unevictable_pgs_stranded 0
thp_fault_alloc 2
thp_fault_fallback 0
thp_collapse_alloc 2
thp_collapse_alloc_failed 0
thp_file_alloc 0
thp_file_mapped 0
thp_split_page 0
thp_split_page_failed 0
thp_deferred_split_page 4
thp_split_pmd 1
thp_split_pud 0
thp_zero_page_alloc 0
thp_zero_page_alloc_failed 0
thp_swpout 0
thp_swpout_fallback 0
balloon_inflate 0
balloon_deflate 0
balloon_migrate 0
swap_ra 0
swap_ra_hit 0
`

func TestVmStatParse(t *testing.T) {
	data, err := parseVMStat([]byte(rawInput))
	if err != nil {
		t.Fatal(err)
	}
	// Check a few values
	assert.Equal(t, uint64(8348560), data.SlabsScanned)
	assert.Equal(t, uint64(0), data.SwapRa)
	assert.Equal(t, uint64(108959), data.WorkingsetActivate)
}
