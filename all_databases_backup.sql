--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE devuser;
ALTER ROLE devuser WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:NXYtN/EYiOH4Xr7+VWJ0dA==$jzqUtWC+ZO7Mr3R/ldB1T0TDOQpz/Ip5DvJwkp0HYL8=:Z3+pl2LxmMSKKiY7q2Aez9eABbv1u/hcKDoQrYvUkyw=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

\connect template1

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Debian 16.9-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Debian 16.9-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

--
-- Database "devdb" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Debian 16.9-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Debian 16.9-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: devdb; Type: DATABASE; Schema: -; Owner: devuser
--

CREATE DATABASE devdb WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE devdb OWNER TO devuser;

\connect devdb

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: goals; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.goals (
    id bigint NOT NULL,
    name character varying(255) NOT NULL,
    target_amount double precision NOT NULL,
    start_amount double precision NOT NULL,
    frequency smallint NOT NULL,
    deadline_date timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL
);


ALTER TABLE public.goals OWNER TO devuser;

--
-- Name: goals_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.goals_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.goals_id_seq OWNER TO devuser;

--
-- Name: goals_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.goals_id_seq OWNED BY public.goals.id;


--
-- Name: prescribed_expanses; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.prescribed_expanses (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    category_id bigint NOT NULL,
    description character varying(255) NOT NULL,
    frequency smallint NOT NULL,
    amount double precision NOT NULL,
    date_time timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL
);


ALTER TABLE public.prescribed_expanses OWNER TO devuser;

--
-- Name: prescribed_expanses_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.prescribed_expanses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.prescribed_expanses_id_seq OWNER TO devuser;

--
-- Name: prescribed_expanses_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.prescribed_expanses_id_seq OWNED BY public.prescribed_expanses.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO devuser;

--
-- Name: transaction_categories; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.transaction_categories (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    name character varying(255) NOT NULL,
    type smallint NOT NULL,
    created_at timestamp without time zone NOT NULL
);


ALTER TABLE public.transaction_categories OWNER TO devuser;

--
-- Name: transaction_categories_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.transaction_categories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transaction_categories_id_seq OWNER TO devuser;

--
-- Name: transaction_categories_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.transaction_categories_id_seq OWNED BY public.transaction_categories.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.transactions (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    category_id bigint NOT NULL,
    goal_id bigint,
    description character varying(255) NOT NULL,
    amount double precision NOT NULL,
    type smallint NOT NULL,
    date_time timestamp without time zone NOT NULL,
    created_at timestamp without time zone NOT NULL,
    prescribed_expanse_id bigint
);


ALTER TABLE public.transactions OWNER TO devuser;

--
-- Name: transactions_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.transactions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transactions_id_seq OWNER TO devuser;

--
-- Name: transactions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.transactions_id_seq OWNED BY public.transactions.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: devuser
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    username text,
    password_hash text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    role smallint NOT NULL
);


ALTER TABLE public.users OWNER TO devuser;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: devuser
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO devuser;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: devuser
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: goals id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.goals ALTER COLUMN id SET DEFAULT nextval('public.goals_id_seq'::regclass);


--
-- Name: prescribed_expanses id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.prescribed_expanses ALTER COLUMN id SET DEFAULT nextval('public.prescribed_expanses_id_seq'::regclass);


--
-- Name: transaction_categories id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transaction_categories ALTER COLUMN id SET DEFAULT nextval('public.transaction_categories_id_seq'::regclass);


--
-- Name: transactions id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transactions ALTER COLUMN id SET DEFAULT nextval('public.transactions_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: goals; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.goals (id, name, target_amount, start_amount, frequency, deadline_date, created_at) FROM stdin;
\.


--
-- Data for Name: prescribed_expanses; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.prescribed_expanses (id, user_id, category_id, description, frequency, amount, date_time, created_at) FROM stdin;
1	1	2	Маме	0	4000	2025-11-15 14:37:54.886427	2025-11-15 14:37:54.886902
2	1	3	Мамина квартира	0	1500	2025-11-15 14:37:54.890813	2025-11-15 14:37:54.890949
3	1	3	Интернет	0	135	2025-11-15 14:37:54.894019	2025-11-15 14:37:54.894155
4	1	8	Еда	0	4500	2025-11-15 14:37:54.897016	2025-11-15 14:37:54.897154
5	1	8	Съёмная квартира	0	2800	2025-11-15 14:37:54.899635	2025-11-15 14:37:54.899767
6	1	8	Комуналка	0	800	2025-11-15 14:37:54.901984	2025-11-15 14:37:54.902123
7	1	8	Бытовая химия	0	600	2025-11-15 14:37:54.904248	2025-11-15 14:37:54.904355
8	1	8	Тренажёрка	0	230	2025-11-15 14:37:54.906392	2025-11-15 14:37:54.906494
9	1	8	Стрижка	0	160	2025-11-15 14:37:54.908506	2025-11-15 14:37:54.908619
10	1	8	Рассрочка	0	900	2025-11-15 14:37:54.910628	2025-11-15 14:37:54.910776
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.schema_migrations (version, dirty) FROM stdin;
20251115000000	f
\.


--
-- Data for Name: transaction_categories; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.transaction_categories (id, user_id, name, type, created_at) FROM stdin;
1	1	Техника	1	2025-11-15 14:37:54.882047
2	1	Продукты	1	2025-11-15 14:37:54.885586
3	1	Другое	1	2025-11-15 14:37:54.888369
4	1	Хавчик вне дома	1	2025-11-15 14:37:54.890109
5	1	Транспорт	1	2025-11-15 14:37:54.891735
6	1	Кафешки	1	2025-11-15 14:37:54.893296
7	1	Комуналка / Квартира	1	2025-11-15 14:37:54.894954
8	1	Бытовая химия	1	2025-11-15 14:37:54.896411
9	1	В дом	1	2025-11-15 14:37:54.898027
10	1	Маме	1	2025-11-15 14:37:54.901364
11	1	Связь телефон / инет	1	2025-11-15 14:37:54.903665
12	1	Медицина	1	2025-11-15 14:37:54.905767
13	1	Отложить	1	2025-11-15 14:37:54.907909
14	1	Подарки	1	2025-11-15 14:37:54.910015
15	1	Одежда	1	2025-11-15 14:37:54.912348
16	1	Красота	1	2025-11-15 14:37:54.916855
17	1	Обучение	1	2025-11-15 14:37:54.918316
18	1	Подписки	1	2025-11-15 14:37:54.919075
19	1	Маме (остальное)	1	2025-11-15 14:37:54.92041
20	1	Развлечения / Отдых	1	2025-11-15 14:37:54.921825
21	1	Машина	1	2025-11-15 14:37:54.923413
22	1	Погрешность	1	2025-11-15 14:37:54.926296
23	1	Документы	1	2025-11-15 14:37:54.927798
24	1	Любимой	1	2025-11-15 14:37:54.929607
25	1	Родителям Карины	1	2025-11-15 14:37:54.932188
26	1	Зарплата	0	2025-11-15 14:39:15.639005
27	1	Мелкий доход	0	2025-11-15 14:39:24.279496
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.transactions (id, user_id, category_id, goal_id, description, amount, type, date, created_at, prescribed_expanse_id) FROM stdin;
1	1	1	\N	Рассрочка	900	1	2025-11-15 14:37:54.881707	2025-11-15 14:37:54.88353	\N
2	1	3	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.888147	2025-11-15 14:37:54.889091	\N
3	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.891618	2025-11-15 14:37:54.892378	\N
4	1	7	\N	Обменял 100 бачей за хату	1630	1	2025-11-15 14:37:54.894848	2025-11-15 14:37:54.895533	\N
5	1	9	\N	Импортированная транзакция	847	1	2025-11-15 14:37:54.8979	2025-11-15 14:37:54.898757	\N
6	1	6	\N	Доставка	289	1	2025-11-15 14:37:54.900437	2025-11-15 14:37:54.900538	\N
7	1	2	\N	Импортированная транзакция	249	1	2025-11-15 14:37:54.902815	2025-11-15 14:37:54.90292	\N
8	1	5	\N	Такси	36	1	2025-11-15 14:37:54.904977	2025-11-15 14:37:54.905081	\N
9	1	5	\N	Такси	23	1	2025-11-15 14:37:54.90711	2025-11-15 14:37:54.907213	\N
10	1	5	\N	Импортированная транзакция	5	1	2025-11-15 14:37:54.909212	2025-11-15 14:37:54.909319	\N
11	1	4	\N	Импортированная транзакция	25	1	2025-11-15 14:37:54.911415	2025-11-15 14:37:54.911527	\N
12	1	6	\N	Кофе	23	1	2025-11-15 14:37:54.913759	2025-11-15 14:37:54.913891	\N
13	1	1	\N	Трещётка и плоскогубцы	316	1	2025-11-15 14:37:54.914813	2025-11-15 14:37:54.914948	\N
14	1	5	\N	Импортированная транзакция	5	1	2025-11-15 14:37:54.915872	2025-11-15 14:37:54.916007	\N
15	1	4	\N	Импортированная транзакция	28	1	2025-11-15 14:37:54.917445	2025-11-15 14:37:54.917551	\N
16	1	18	\N	VDS'ка	288	1	2025-11-15 14:37:54.918947	2025-11-15 14:37:54.919698	\N
17	1	6	\N	Импортированная транзакция	23	1	2025-11-15 14:37:54.920982	2025-11-15 14:37:54.921096	\N
19	1	21	\N	Буксировка	200	1	2025-11-15 14:37:54.923305	2025-11-15 14:37:54.923969	\N
20	1	21	\N	Бенз	600	1	2025-11-15 14:37:54.924682	2025-11-15 14:37:54.924789	\N
21	1	2	\N	Импортированная транзакция	117	1	2025-11-15 14:37:54.92551	2025-11-15 14:37:54.925617	\N
22	1	12	\N	Импортированная транзакция	38	1	2025-11-15 14:37:54.926861	2025-11-15 14:37:54.92698	\N
23	1	5	\N	Импортированная транзакция	5	1	2025-11-15 14:37:54.928558	2025-11-15 14:37:54.928714	\N
24	1	12	\N	Лекарства	243	1	2025-11-15 14:37:54.930251	2025-11-15 14:37:54.930385	\N
25	1	4	\N	Импортированная транзакция	21	1	2025-11-15 14:37:54.931165	2025-11-15 14:37:54.9313	\N
26	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.932811	2025-11-15 14:37:54.932967	\N
27	1	6	\N	Импортированная транзакция	39	1	2025-11-15 14:37:54.933587	2025-11-15 14:37:54.933721	\N
28	1	5	\N	Импортированная транзакция	5	1	2025-11-15 14:37:54.934419	2025-11-15 14:37:54.934537	\N
29	1	2	\N	Импортированная транзакция	78	1	2025-11-15 14:37:54.935156	2025-11-15 14:37:54.935273	\N
30	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.935908	2025-11-15 14:37:54.936027	\N
31	1	9	\N	Лавашик	50	1	2025-11-15 14:37:54.936608	2025-11-15 14:37:54.936779	\N
32	1	6	\N	Импортированная транзакция	168	1	2025-11-15 14:37:54.937442	2025-11-15 14:37:54.937551	\N
33	1	2	\N	Импортированная транзакция	40	1	2025-11-15 14:37:54.938256	2025-11-15 14:37:54.938417	\N
34	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.939125	2025-11-15 14:37:54.939244	\N
35	1	2	\N	Импортированная транзакция	16	1	2025-11-15 14:37:54.939818	2025-11-15 14:37:54.93994	\N
37	1	11	\N	Импортированная транзакция	25	1	2025-11-15 14:37:54.941437	2025-11-15 14:37:54.941565	\N
38	1	20	\N	Кинчик	58	1	2025-11-15 14:37:54.942216	2025-11-15 14:37:54.942348	\N
39	1	20	\N	Попкорн	39	1	2025-11-15 14:37:54.943023	2025-11-15 14:37:54.943158	\N
40	1	5	\N	Импортированная транзакция	35	1	2025-11-15 14:37:54.943752	2025-11-15 14:37:54.943867	\N
41	1	7	\N	20 бачей	327	1	2025-11-15 14:37:54.944456	2025-11-15 14:37:54.944579	\N
42	1	8	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.945176	2025-11-15 14:37:54.945294	\N
43	1	4	\N	Импортированная транзакция	23	1	2025-11-15 14:37:54.945935	2025-11-15 14:37:54.946077	\N
44	1	2	\N	Импортированная транзакция	41	1	2025-11-15 14:37:54.946789	2025-11-15 14:37:54.946918	\N
45	1	2	\N	Импортированная транзакция	138	1	2025-11-15 14:37:54.947526	2025-11-15 14:37:54.947652	\N
46	1	13	\N	Импортированная транзакция	4500	1	2025-11-15 14:37:54.948242	2025-11-15 14:37:54.948365	\N
47	1	11	\N	Безлимит	100	1	2025-11-15 14:37:54.948941	2025-11-15 14:37:54.949068	\N
48	1	5	\N	Такси до бабушки	31	1	2025-11-15 14:37:54.949709	2025-11-15 14:37:54.949831	\N
49	1	7	\N	50 бачей	830	1	2025-11-15 14:37:54.950427	2025-11-15 14:37:54.950554	\N
50	1	2	\N	Импортированная транзакция	124	1	2025-11-15 14:37:54.951166	2025-11-15 14:37:54.951303	\N
51	1	5	\N	Импортированная транзакция	5	1	2025-11-15 14:37:54.951893	2025-11-15 14:37:54.952006	\N
52	1	4	\N	Импортированная транзакция	40	1	2025-11-15 14:37:54.95264	2025-11-15 14:37:54.952752	\N
53	1	5	\N	Маме такси	32	1	2025-11-15 14:37:54.953305	2025-11-15 14:37:54.953422	\N
54	1	4	\N	Импортированная транзакция	18	1	2025-11-15 14:37:54.953993	2025-11-15 14:37:54.954111	\N
55	1	25	\N	Импортированная транзакция	125	1	2025-11-15 14:37:54.95467	2025-11-15 14:37:54.954784	\N
56	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.955356	2025-11-15 14:37:54.955468	\N
57	1	2	\N	Мяса	71	1	2025-11-15 14:37:54.956031	2025-11-15 14:37:54.956208	\N
58	1	2	\N	Импортированная транзакция	84	1	2025-11-15 14:37:54.9568	2025-11-15 14:37:54.957011	\N
59	1	5	\N	Импортированная транзакция	35	1	2025-11-15 14:37:54.957572	2025-11-15 14:37:54.957851	\N
60	1	2	\N	Импортированная транзакция	60	1	2025-11-15 14:37:54.958637	2025-11-15 14:37:54.958772	\N
61	1	15	\N	Куртка	1150	1	2025-11-15 14:37:54.959373	2025-11-15 14:37:54.959489	\N
62	1	4	\N	Импортированная транзакция	38	1	2025-11-15 14:37:54.960189	2025-11-15 14:37:54.960316	\N
63	1	2	\N	Импортированная транзакция	99	1	2025-11-15 14:37:54.960924	2025-11-15 14:37:54.961041	\N
64	1	5	\N	Импортированная транзакция	5	1	2025-11-15 14:37:54.961642	2025-11-15 14:37:54.961786	\N
65	1	16	\N	Стрижка	160	1	2025-11-15 14:37:54.962842	2025-11-15 14:37:54.962998	\N
66	1	5	\N	Маме такси	34	1	2025-11-15 14:37:54.963689	2025-11-15 14:37:54.96387	\N
67	1	5	\N	Такси от бабушки	29	1	2025-11-15 14:37:54.96454	2025-11-15 14:37:54.964662	\N
68	1	2	\N	Импортированная транзакция	6	1	2025-11-15 14:37:54.965252	2025-11-15 14:37:54.965377	\N
69	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.966001	2025-11-15 14:37:54.966127	\N
70	1	25	\N	Импортированная транзакция	400	1	2025-11-15 14:37:54.966732	2025-11-15 14:37:54.966846	\N
71	1	2	\N	Импортированная транзакция	52	1	2025-11-15 14:37:54.967436	2025-11-15 14:37:54.967552	\N
72	1	2	\N	Импортированная транзакция	8	1	2025-11-15 14:37:54.968161	2025-11-15 14:37:54.96827	\N
73	1	19	\N	Продукты	159	1	2025-11-15 14:37:54.968834	2025-11-15 14:37:54.968952	\N
74	1	4	\N	Импортированная транзакция	20	1	2025-11-15 14:37:54.969526	2025-11-15 14:37:54.969638	\N
75	1	5	\N	Импортированная транзакция	14	1	2025-11-15 14:37:54.9702	2025-11-15 14:37:54.970316	\N
77	1	14	\N	Бабушке	300	1	2025-11-15 14:37:54.971634	2025-11-15 14:37:54.97175	\N
78	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.972339	2025-11-15 14:37:54.972461	\N
79	1	6	\N	Импортированная транзакция	18	1	2025-11-15 14:37:54.973234	2025-11-15 14:37:54.973484	\N
80	1	2	\N	Импортированная транзакция	56	1	2025-11-15 14:37:54.974187	2025-11-15 14:37:54.974321	\N
81	1	5	\N	Такси 6 числа	45	1	2025-11-15 14:37:54.97493	2025-11-15 14:37:54.975064	\N
82	1	5	\N	Импортированная транзакция	15	1	2025-11-15 14:37:54.975712	2025-11-15 14:37:54.975839	\N
83	1	4	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.976448	2025-11-15 14:37:54.976564	\N
84	1	11	\N	4г	55	1	2025-11-15 14:37:54.977155	2025-11-15 14:37:54.977273	\N
85	1	5	\N	Импортированная транзакция	10	1	2025-11-15 14:37:54.977887	2025-11-15 14:37:54.978002	\N
86	1	16	\N	ноготочки	180	1	2025-11-15 14:37:54.978592	2025-11-15 14:37:54.978702	\N
87	1	20	\N	пряжа	89	1	2025-11-15 14:37:54.979327	2025-11-15 14:37:54.979456	\N
88	1	26	\N		24330	0	2025-11-15 10:39:00	2025-11-15 14:40:13.546267	\N
89	1	27	\N	Были	580	0	2025-11-15 10:40:00	2025-11-15 14:40:30.799525	\N
90	1	27	\N		300	0	2025-11-15 10:40:00	2025-11-15 14:41:06.24088	\N
91	1	26	\N		1010	0	2025-11-15 10:41:00	2025-11-15 14:41:16.610909	\N
92	1	27	\N		161	0	2025-11-15 10:41:00	2025-11-15 14:41:36.630519	\N
93	1	3	\N	Интернет	135	1	2025-11-15 14:42:12.117431	2025-11-15 14:42:12.117431	3
98	1	2	\N	Маме	3000	1	2025-11-15 13:35:00	2025-11-15 15:35:40.786126	1
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: devuser
--

COPY public.users (id, username, password_hash, created_at, role) FROM stdin;
1	test	$2a$10$ml8lmrkUXzQ62KdHo2s3Ge455YNgtjJDGGOzCotW3m0kN4RqSW79G	2025-11-15 14:37:12.080581	0
\.


--
-- Name: goals_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.goals_id_seq', 1, false);


--
-- Name: prescribed_expanses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.prescribed_expanses_id_seq', 11, true);


--
-- Name: transaction_categories_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.transaction_categories_id_seq', 27, true);


--
-- Name: transactions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.transactions_id_seq', 98, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devuser
--

SELECT pg_catalog.setval('public.users_id_seq', 1, true);


--
-- Name: goals goals_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.goals
    ADD CONSTRAINT goals_pkey PRIMARY KEY (id);


--
-- Name: prescribed_expanses prescribed_expanses_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.prescribed_expanses
    ADD CONSTRAINT prescribed_expanses_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: transaction_categories transaction_categories_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transaction_categories
    ADD CONSTRAINT transaction_categories_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: prescribed_expanses prescribed_expanses_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.prescribed_expanses
    ADD CONSTRAINT prescribed_expanses_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.transaction_categories(id);


--
-- Name: prescribed_expanses prescribed_expanses_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.prescribed_expanses
    ADD CONSTRAINT prescribed_expanses_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: transaction_categories transaction_categories_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transaction_categories
    ADD CONSTRAINT transaction_categories_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: transactions transactions_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.transaction_categories(id);


--
-- Name: transactions transactions_goal_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_goal_id_fkey FOREIGN KEY (goal_id) REFERENCES public.goals(id);


--
-- Name: transactions transactions_prescribed_expanse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_prescribed_expanse_id_fkey FOREIGN KEY (prescribed_expanse_id) REFERENCES public.prescribed_expanses(id);


--
-- Name: transactions transactions_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: devuser
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

\connect postgres

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.9 (Debian 16.9-1.pgdg120+1)
-- Dumped by pg_dump version 16.9 (Debian 16.9-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

