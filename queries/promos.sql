SELECT 
    p.id, 
    p.title, 
    p.status, 
    p.sweepz_entry as sweepz_per_entry, 
    count(pe.id) as entries, 
    count(distinct pe.playfab_id) as unique_entries, 
    count(pe.id) * p.sweepz_entry as total_sweepz, 
    p.start_at, 
    p.end_at
FROM promo_entries pe
LEFT JOIN promos p on pe.promo_id = p.id
WHERE p.start_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY pe.promo_id
ORDER BY p.id 