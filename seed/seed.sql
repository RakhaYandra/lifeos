-- Seed P6 fiktif (bukan data asli). Password semua: Rahasia123
-- Persona: pekerja Jakarta, sewa kos, gajian tiap tanggal 1.
-- Relasi: goals annual->monthly->projects->tasks, trx->budgets.
INSERT INTO users (email, password_hash) VALUES
('aku@lifeos.local', '$2a$10$GMr6NAxuAxPKutSu2J5MUO3DG.3iqRcYv5Gu9NQNMnTNRErbsKqMu');

INSERT INTO settings (user_id, active_year, currency, budget_warn_pct, goal_warn_pct) VALUES
(1, 2026, 'IDR', 80, 70);

-- Goals: 4 annual + 3 monthly (parent)
INSERT INTO goals (user_id,level,parent_id,life_area_id,title,metric,target_value,current_value,status,target_date) VALUES
(1,'annual',NULL,2,'Dana darurat 20jt','IDR',20000000,8000000,'active','2026-12-31'),
(1,'annual',NULL,3,'Lari 10K tanpa berhenti','km',10,6.5,'on_track','2026-11-30'),
(1,'annual',NULL,4,'Go intermediate','modul',12,5,'active','2026-12-31'),
(1,'annual',NULL,5,'Baca 24 buku','buku',24,14,'on_track','2026-12-31'),
(1,'monthly',1,2,'Nabung September','IDR',2000000,800000,'active','2026-09-30'),
(1,'monthly',2,3,'Lari 30km September','km',30,18,'active','2026-09-30'),
(1,'monthly',3,4,' modul Go pekan 2','modul',2,1,'active','2026-09-30');

-- Projects (goal bulanan where relevant)
INSERT INTO projects (user_id,name,life_area_id,goal_id,status,start_date,target_date,notes) VALUES
(1,'Renovasi kamar kos',8,NULL,'active','2026-09-01','2026-10-15','cat + rak'),
(1,'LifeOS MVP',1,7,'active','2026-09-10','2026-09-30','Go+SQLite+React'),
(1,'Persiapan lari 10K',3,6,'active','2026-09-01','2026-09-30','jadwal latihan'),
(1,'Beresi baca Q3',5,NULL,'planning','2026-09-15','2026-09-30','3 buku'),
(1,'Laporan keuangan Q3',2,5,'planning','2026-09-20','2026-09-30','rekap tabungan');

-- Tasks: 25 (campur status; 2 overdue: due 09-09 & 09-10 non-completed; 3 due 09-14)
INSERT INTO tasks (user_id,project_id,goal_id,title,status,priority,due_date,completed_at,life_area_id) VALUES
(1,1,NULL,'Beli cat 5kg + kuas','completed','high','2026-09-02','2026-09-02',8),
(1,1,NULL,'Cat dinding kamar','in_progress','high','2026-09-10',NULL,8),
(1,1,NULL,'Rakit rak buku','not_started','medium','2026-09-14',NULL,8),
(1,1,NULL,'Buang kardus bekas','inbox','low',NULL,NULL,8),
(1,2,7,'Scaffold API + auth','completed','high','2026-09-11','2026-09-11',1),
(1,2,7,'CRUD tasks + overdue logic','completed','high','2026-09-12','2026-09-12',1),
(1,2,7,'Streak habits + heatmap','in_progress','high','2026-09-14',NULL,1),
(1,2,7,'Finance lite + budget aktual','waiting','medium','2026-09-16',NULL,1),
(1,2,7,'Tulis README 5 menit','not_started','low','2026-09-20',NULL,1),
(1,3,6,'Lari interval 5K','completed','medium','2026-09-08','2026-09-08',3),
(1,3,6,'Long run 8K','in_progress','high','2026-09-11',NULL,3),
(1,3,6,'Beli sepatu baru','completed','medium','2026-09-05','2026-09-05',3),
(1,3,6,'Long run 10K','not_started','high','2026-09-14',NULL,3),
(1,4,NULL,'Selesaikan Atomic Habits','in_progress','medium','2026-09-14',NULL,5),
(1,4,NULL,'Resensi 1 halaman','inbox','low',NULL,NULL,5),
(1,NULL,5,'Transfer tabungan 800rb','completed','high','2026-09-01','2026-09-01',2),
(1,NULL,5,'Review budget Agustus','completed','medium','2026-09-03','2026-09-03',2),
(1,NULL,NULL,'Bayar kos Oktober','not_started','critical','2026-09-28',NULL,8),
(1,NULL,NULL,'Servis motor','waiting','medium','2026-09-09',NULL,8),
(1,NULL,NULL,'Telemedicine gigi','cancelled','low','2026-09-06',NULL,3),
(1,NULL,NULL,'Update CV + porto','in_progress','medium','2026-09-18',NULL,1),
(1,NULL,NULL,'Belajar generics Go','completed','medium','2026-09-07','2026-09-07',4),
(1,NULL,NULL,'Telepon ibu','completed','medium','2026-09-13','2026-09-13',7),
(1,NULL,NULL,'Bereskan meja kerja','inbox','low',NULL,NULL,8),
(1,NULL,NULL,'Donor darah','not_started','low','2026-09-21',NULL,3);

-- Habits: 3 daily + 2 weekly
INSERT INTO habits (user_id,name,life_area_id,frequency,target_per_week,start_date,active) VALUES
(1,'Olahraga 20 menit',3,'daily',1,'2026-08-16',1),
(1,'Baca 30 halaman',5,'daily',1,'2026-08-16',1),
(1,'Meditasi 10 menit',5,'daily',1,'2026-09-01',1),
(1,'Lari',3,'weekly',3,'2026-08-16',1),
(1,'Review mingguan',1,'weekly',1,'2026-08-16',1);

-- Habit logs: 80 baris 30 hari (streak h1=14, h2 bolong, h4 3x/pekan)
INSERT INTO habit_logs (habit_id,date,done) VALUES
(1,'2026-08-16',1),
(1,'2026-08-17',1),
(1,'2026-08-18',1),
(1,'2026-08-19',1),
(1,'2026-08-21',1),
(1,'2026-08-22',1),
(1,'2026-08-23',1),
(1,'2026-08-24',1),
(1,'2026-08-26',1),
(1,'2026-08-27',1),
(1,'2026-08-28',1),
(1,'2026-08-29',1),
(1,'2026-08-30',1),
(1,'2026-08-31',1),
(1,'2026-09-01',1),
(1,'2026-09-02',1),
(1,'2026-09-03',1),
(1,'2026-09-04',1),
(1,'2026-09-05',1),
(1,'2026-09-06',1),
(1,'2026-09-07',1),
(1,'2026-09-08',1),
(1,'2026-09-09',1),
(1,'2026-09-10',1),
(1,'2026-09-11',1),
(1,'2026-09-12',1),
(1,'2026-09-13',1),
(1,'2026-09-14',1),
(2,'2026-08-16',1),
(2,'2026-08-17',1),
(2,'2026-08-19',1),
(2,'2026-08-20',1),
(2,'2026-08-22',1),
(2,'2026-08-23',1),
(2,'2026-08-25',1),
(2,'2026-08-26',1),
(2,'2026-08-28',1),
(2,'2026-08-29',1),
(2,'2026-08-31',1),
(2,'2026-09-01',1),
(2,'2026-09-02',1),
(2,'2026-09-04',1),
(2,'2026-09-05',1),
(2,'2026-09-07',1),
(2,'2026-09-08',1),
(2,'2026-09-10',1),
(2,'2026-09-11',1),
(2,'2026-09-12',0),
(2,'2026-09-13',1),
(2,'2026-09-14',1),
(3,'2026-09-01',1),
(3,'2026-09-02',1),
(3,'2026-09-03',1),
(3,'2026-09-04',1),
(3,'2026-09-06',1),
(3,'2026-09-07',1),
(3,'2026-09-08',1),
(3,'2026-09-10',1),
(3,'2026-09-11',1),
(3,'2026-09-12',1),
(3,'2026-09-13',1),
(3,'2026-09-14',1),
(4,'2026-08-16',1),
(4,'2026-08-18',1),
(4,'2026-08-20',1),
(4,'2026-08-23',1),
(4,'2026-08-25',1),
(4,'2026-08-27',1),
(4,'2026-08-30',1),
(4,'2026-09-01',1),
(4,'2026-09-03',1),
(4,'2026-09-06',1),
(4,'2026-09-08',1),
(4,'2026-09-10',1),
(4,'2026-09-13',1),
(5,'2026-08-16',1),
(5,'2026-08-23',1),
(5,'2026-08-30',1),
(5,'2026-09-06',1),
(5,'2026-09-13',1);

-- Transactions: 32 (IDR; gaji+freelance+kos+makan dst)
INSERT INTO transactions (user_id,date,type,category,description,amount,account,recurring) VALUES
(1,'2026-09-01','income','Gaji','Gaji September',8500000,'BCA',1),
(1,'2026-09-05','income','Freelance','Landing page UMKM',1500000,'BCA',0),
(1,'2026-09-01','expense','Hunian','Kos Oktober dibayar awal',1000000,'BCA',1),
(1,'2026-09-01','expense','Makan','Warteg',25000,'Tunai',0),
(1,'2026-09-02','expense','Makan','Padang',35000,'Tunai',0),
(1,'2026-09-02','expense','Transport','TransJakarta',7000,'E-money',0),
(1,'2026-09-03','expense','Makan','Soto + es teh',28000,'Tunai',0),
(1,'2026-09-03','expense','Belanja','Sabun + sampo',65000,'BCA',0),
(1,'2026-09-04','expense','Makan','Nasi goreng',22000,'Tunai',0),
(1,'2026-09-04','expense','Hiburan','Bioskop',75000,'BCA',0),
(1,'2026-09-05','expense','Transport','Bensin motor',50000,'Tunai',0),
(1,'2026-09-05','expense','Makan','Ayam geprek',18000,'Tunai',0),
(1,'2026-09-06','expense','Makan','Bakso',20000,'Tunai',0),
(1,'2026-09-06','expense','Kesehatan','Vitamin',85000,'BCA',0),
(1,'2026-09-07','expense','Makan','Warung Bu Tini',30000,'Tunai',0),
(1,'2026-09-07','expense','Internet','Paket data 30GB',100000,'BCA',1),
(1,'2026-09-08','expense','Makan','Sate padang',32000,'Tunai',0),
(1,'2026-09-08','expense','Transport','Ojol ke klien',45000,'BCA',0),
(1,'2026-09-09','expense','Makan','Mie ayam',18000,'Tunai',0),
(1,'2026-09-09','expense','Edukasi','Buku Atomic Habits',99000,'BCA',0),
(1,'2026-09-10','expense','Makan','Pecel lele',25000,'Tunai',0),
(1,'2026-09-10','expense','Hunian','Cat + kuas',185000,'BCA',0),
(1,'2026-09-11','expense','Makan','Gado-gado',20000,'Tunai',0),
(1,'2026-09-11','expense','Transport','Parkir + tol',35000,'E-money',0),
(1,'2026-09-12','expense','Makan','Nasi padang',30000,'Tunai',0),
(1,'2026-09-12','expense','Hiburan','Spotify',54999,'BCA',1),
(1,'2026-09-13','expense','Makan','Sop iga',40000,'Tunai',0),
(1,'2026-09-13','expense','Keluarga','Transfer ibu',500000,'BCA',0),
(1,'2026-09-14','expense','Makan','Bubur ayam',15000,'Tunai',0),
(1,'2026-09-14','expense','Transport','TransJakarta',7000,'E-money',0),
(1,'2026-08-28','expense','Makan','Kantin (Agustus)',22000,'Tunai',0),
(1,'2026-08-15','income','Gaji','Gaji Agustus (riwayat)',8500000,'BCA',0);

-- Budgets Sep 2026 (aktual terhitung dari trx expense)
INSERT INTO budgets (user_id,year,month,category,amount) VALUES
(1,2026,9,'Makan',1500000),
(1,2026,9,'Transport',500000),
(1,2026,9,'Hiburan',300000),
(1,2026,9,'Hunian',1200000);

-- Subscriptions (1 renewal dekat: iCloud 18 Sep)
INSERT INTO subscriptions (user_id,service,category,cost,frequency,next_billing,payment_method,auto_renew,active,notes) VALUES
(1,'Spotify','Hiburan',54999,'monthly','2026-09-20','BCA',1,1,''),
(1,'Netflix','Hiburan',186000,'monthly','2026-09-25','BCA',1,1,'paket mobile'),
(1,'iCloud 50GB','Digital',15000,'monthly','2026-09-18','BCA',1,1,''),
(1,'Domain rakha.dev','Digital',150000,'yearly','2027-01-10','BCA',1,1,'');

-- Health logs 14 hari
INSERT INTO health_logs (user_id,date,weight,sleep_hours,water_liters,energy,mood,notes) VALUES
(1,'2026-09-01',70.7,6.0,2.0,3,4,''),
(1,'2026-09-02',70.8,6.5,2.0,4,5,''),
(1,'2026-09-03',70.7,7.0,2.0,5,3,''),
(1,'2026-09-04',70.7,7.5,2.0,3,4,''),
(1,'2026-09-05',70.6,6.0,2.0,4,5,''),
(1,'2026-09-06',70.7,6.5,2.0,5,3,''),
(1,'2026-09-07',70.6,7.0,2.0,3,4,''),
(1,'2026-09-08',70.6,7.5,2.0,4,5,''),
(1,'2026-09-09',70.5,6.0,2.0,5,3,''),
(1,'2026-09-10',70.6,6.5,2.0,3,4,''),
(1,'2026-09-11',70.5,7.0,2.0,4,5,''),
(1,'2026-09-12',70.5,7.5,2.0,5,3,''),
(1,'2026-09-13',70.4,6.0,2.0,3,4,''),
(1,'2026-09-14',70.5,6.5,2.0,4,5,'');

-- Workouts 8x
INSERT INTO workouts (user_id,date,type,duration_min,intensity,calories,notes) VALUES
(1,'2026-09-01','Lari',25,'medium',250,''),
(1,'2026-09-03','Lari',30,'high',320,''),
(1,'2026-09-05','Bodyweight',20,'medium',180,''),
(1,'2026-09-07','Lari',45,'high',480,''),
(1,'2026-09-08','Jalan',30,'low',120,''),
(1,'2026-09-10','Lari',35,'high',380,''),
(1,'2026-09-12','Bodyweight',20,'medium',190,''),
(1,'2026-09-14','Lari',30,'medium',300,'');

-- Learning + reading
INSERT INTO learning_entries (user_id,topic,type,provider,related_skill,progress,hours,status,notes) VALUES
(1,'Go generics + testing','course','YouTube - Golang ID','Go',40,5,'active',''),
(1,'React hooks deep dive','tutorial','Web Dev Simplified','React',70,8,'active','');
INSERT INTO reading_entries (user_id,title,type,author,status,rating,takeaways) VALUES
(1,'Atomic Habits','book','James Clear','finished',5,'lingkungan > motivasi'),
(1,'Clean Code','book','R. C. Martin','reading',NULL,'nama jelas, fungsi kecil'),
(1,'Go by Example','docs','gobyexample.com','reading',NULL,'patterns konkurensi');

-- Review pekan 8 Sep + reminders
INSERT INTO reviews (user_id,week_start,stats,wins,challenges,lessons,next_focus) VALUES
(1,'2026-09-07','{}','API P1-P3 selesai','long run kelewat 1x','tidur <7 jam bikin lesu','P4 + long run 10K');
INSERT INTO reminders (user_id,title,date,recurrence,notes) VALUES
(1,'Bayar kos Oktober','2026-10-01','monthly',''),
(1,'Ulang tahun ibu','2026-11-20','yearly',''),
(1,'Pajak motor','2027-02-15','yearly',''),
(1,'Kontrol gigi','2026-09-22','none',''),
(1,'Long run 10K','2026-09-14','weekly','');

-- Fase 2B: monthly + yearly reviews (stats dihitung API saat create)
INSERT INTO monthly_reviews (user_id,period,stats,wins,challenges,lessons,next_focus) VALUES
(1,'2026-08','{}','Gaji aman, lari mulai rutin','kurang catat trx','catat harian lebih gampang','nabung 2jt September'),
(1,'2026-09','{}','API LifeOS P1-P6 selesai','long run kelewat 1x','tidur <7 jam bikin lesu','long run 10K + P7');
INSERT INTO yearly_reviews (user_id,period,stats,wins,challenges,lessons,next_focus,achievements,next_year) VALUES
(1,'2026','{}','LifeOS MVP live','konsistensi lari','sistem > motivasi','marathon 2027','API + web LifeOS, dana darurat 40%','half marathon + tabungan 20jt');

-- Fase 2C: travel contoh + decision contoh
INSERT INTO trips (user_id,name,destination,start_date,end_date,budget,status,notes) VALUES
(1,'Bandung 2 hari','Bandung','2026-10-10','2026-10-11',1500000,'planning','long weekend');
INSERT INTO itinerary_items (trip_id,date,time,activity,location,cost,booked) VALUES
(1,'2026-10-10','08:00','Kereta Gambir-Kiaracondang','Gambir',150000,1),
(1,'2026-10-10','13:00','Kawah Putih','Ciwidey',50000,0),
(1,'2026-10-11','09:00','Braga + kopi','Braga',80000,0),
(1,'2026-10-11','16:00','Kereta pulang','Kiaracondang',150000,1);
INSERT INTO packing_items (trip_id,category,item,qty,packed) VALUES
(1,'Pakaian','Jaket',1,0),
(1,'Pakaian','Baju ganti',2,0),
(1,'Elektronik','Charger + powerbank',1,1),
(1,'Dokumen','KTP',1,1);
INSERT INTO decisions (user_id,title,notes) VALUES
(1,'Pilih laptop baru','ganti 2026'),
(1,'Liburan akhir tahun','');
INSERT INTO decision_options (decision_id,name) VALUES
(1,'ThinkPad X1'),(1,'MacBook Air'),(2,'Bandung'),(2,'Yogya');
INSERT INTO decision_marks (option_id,criterion,weight,score) VALUES
(1,'Harga',3,7),(1,'Bobot',2,9),(1,'Linux',3,10),
(2,'Harga',3,5),(2,'Bobot',2,10),(2,'Linux',3,4),
(3,'Biaya',3,9),(3,'Jarak',2,9),
(4,'Biaya',3,7),(4,'Jarak',2,6);

-- Fase 2A: quarterly goals (cascade annual->quarterly) + reparent monthly
INSERT INTO goals (user_id,level,parent_id,life_area_id,title,metric,target_value,current_value,status,target_date) VALUES
(1,'quarterly',1,2,'Q3 nabung 6jt','IDR',6000000,2800000,'active','2026-09-30'),
(1,'quarterly',2,3,'Q3 base 60km','km',60,42,'on_track','2026-09-30'),
(1,'quarterly',3,4,'Q3 Go 6 modul','modul',6,3,'active','2026-09-30');
UPDATE goals SET parent_id=8 WHERE id=5;
UPDATE goals SET parent_id=9 WHERE id=6;
UPDATE goals SET parent_id=10 WHERE id=7;

-- Milestones (1 overdue contoh: 2026-09-09)
INSERT INTO milestones (user_id,goal_id,project_id,title,target_date,status,completed_at,notes) VALUES
(1,9,3,'Long run 8K tembus','2026-09-09','pending',NULL,''),
(1,9,3,'Long run 10K','2026-09-14','pending',NULL,''),
(1,8,NULL,'Transfer 800rb pekan 2','2026-09-15','pending',NULL,''),
(1,10,2,'Heatmap 30 hari live','2026-09-14','completed','2026-09-14',''),
(1,NULL,1,'Kamar selesai dicat','2026-09-20','pending',NULL,'');
