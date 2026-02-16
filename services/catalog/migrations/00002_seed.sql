-- +goose Up
-- +goose StatementBegin
INSERT INTO products (id, name, description, category, base_price, is_available, image_url) VALUES
-- Классика (Category 0)
('019c53e1-4f90-7373-8d00-6472a7dc310c', 'Маргарита', 'Классическая итальянская пицца с томатным соусом Сан-Марцано, свежей моцареллой буффало, базиликом и оливковым маслом extra virgin.', 0, 590.00, TRUE, 'https://images.unsplash.com/photo-1574071318508-1cdbab80d002?auto=format&fit=crop&w=800&q=80'),
('019c53e1-4f90-7373-8d00-6472a7dc310d', 'Пепперони Фреш', 'Двойная порция острых колбасок пепперони на подушке из нежной моцареллы и фирменного томатного соуса.', 0, 650.00, TRUE, 'https://images.unsplash.com/photo-1628840042765-356cda07504e?auto=format&fit=crop&w=800&q=80'),
('019c53e1-4f90-7373-8d00-6472a7dc310f', 'Четыре сыра', 'Сливочное сочетание Моцареллы, Горгонзолы, Пармезана и Эмменталя. Настоящая мечта сыромана.', 0, 750.00, TRUE, 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=800&q=80'),

-- Премиум (Category 1)
('019c53e1-4f90-7373-8d00-6472a7dc3111', 'Трюфельная с грибами', 'Белая пицца с обжаренными лесными грибами, трюфельным маслом, тимьяном и моцареллой.', 1, 890.00, TRUE, 'https://images.unsplash.com/photo-1513104890138-7c749659a591?auto=format&fit=crop&w=800&q=80'),
('019c53e1-4f90-7373-8d00-6472a7dc3110', 'Цыпленок Барбекю', 'Гриль-курица, красный лук и кинза на базе соуса барбекю, под слоем тягучей моцареллы.', 1, 790.00, TRUE, 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?auto=format&fit=crop&w=800&q=80'),

-- Вегетарианская (Category 2)
('019c53e1-4f90-7373-8d00-6472a7dc3112', 'Овощной Микс', 'Сладкий перец, маслины, красный лук, томаты и свежий шпинат на томатной основе.', 2, 620.00, TRUE, 'https://images.unsplash.com/photo-1571407970349-bc81e7e96d47?auto=format&fit=crop&w=800&q=80'),

-- Острая (Category 3)
('019c53e1-4f90-7373-8d00-6472a7dc3113', 'Мясная Халапеньо', 'Бекон, ветчина, говядина и очень много перца халапеньо для любителей поострее.', 3, 720.00, TRUE, 'https://images.unsplash.com/photo-1585238342024-78d387f4a707?auto=format&fit=crop&w=800&q=80'),

-- Напитки (Category 4)
('019c53e1-4f90-7373-8d00-6472a7dc3117', 'Кола Крафт', 'Ремесленная кола на натуральных травах и специях. 0.5л', 4, 150.00, TRUE, 'https://images.unsplash.com/photo-1554866585-cd94860890b7?auto=format&fit=crop&w=800&q=80'),
('019c53e1-4f90-7373-8d00-6472a7dc3118', 'Домашний Лимонад', 'Свежевыжатые лимоны, мята и капля меда. 0.5л', 4, 180.00, TRUE, 'https://images.unsplash.com/photo-1513558161293-cdaf765ed2fd?auto=format&fit=crop&w=800&q=80'),
('019c53e1-4f90-7373-8d00-6472a7dc3119', 'Апельсиновый сок', '100% натуральный сок прямого отжима. 0.3л', 4, 220.00, TRUE, 'https://images.unsplash.com/photo-1613478223719-2ab802602423?auto=format&fit=crop&w=800&q=80'),

-- Десерты (Category 5)
('019c53e1-4f90-7373-8d00-6472a7dc3115', 'Чизкейк Нью-Йорк', 'Классический сливочный чизкейк на песочной основе с ягодным конфитюром.', 5, 320.00, TRUE, 'https://images.unsplash.com/photo-1524351199678-941a58a3df50?auto=format&fit=crop&w=800&q=80'),
('019c53e1-4f90-7373-8d00-6472a7dc3116', 'Шоколадный Брауни', 'Теплый, тягучий брауни с кусочками грецкого ореха и бельгийским шоколадом.', 5, 280.00, TRUE, 'https://images.unsplash.com/photo-1515037893149-de7f840978e2?auto=format&fit=crop&w=800&q=80');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE id LIKE '019c53e1-4f90-7373-8d00-6472a7dc31%';
-- +goose StatementEnd
