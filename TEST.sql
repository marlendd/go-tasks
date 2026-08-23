--// Нужно описать модель библиотеки. Есть 3 сущности: "Автор", "Книга", "Читатель"
--// Физически книга только одна и может быть только у одного читателя. Нужно составить
--// таблицы для библиотеки так что бы это учесть

CREATE TABLE authors (
	id int PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE readers (
	id int PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE books (
	id int PRIMARY KEY,
	title TEXT NOT NULL,

	reader_id int REFERENCES readers(id) ON DELETE SET NULL
);

CREATE TABLE books_authors (
    book_id int NOT NULL REFERENCES books (id) ON DELETE CASCADE,
    author_id int NOT NULL REFERENCES authors (id) ON DELETE CASCADE,

    PRIMARY KEY(book_id, author_id)
);

-- выбрать названия всех книг которые на руках

SELECT title 
FROM books
WHERE reader_id IS NOT NULL;

-- Написать запрос - выбрать названия всех книг в библиотеке у которых больше 3 авторов
SELECT b.title
FROM books b
    JOIN books_authors ba ON ba.book_id = b.id
GROUP BY b.id, b.title
HAVING COUNT(*) > 3

-- Написать запрос - выбрать имена топ 3 читаемых авторов на данный момент

SELECT a.name 
FROM authors a 
    JOIN books_authors ba ON a.id = ba.author_id
    JOIN books b ON ba.book_id = b.id
GROUP BY ba.author_id, a.name
WHERE reader_id IS NOT NULL
ORDER BY COUNT(*) DESC
LIMIT 3