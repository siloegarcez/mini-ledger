begin;

create table if not exists transactions (
    id bigint generated always as identity primary key,
    account_id bigint not null references accounts (id) on delete restrict,
    operation_type_id bigint not null references operations_types (operation_type_id) on delete restrict,
    amount bigint not null,
    currency text not null,
    scale smallint not null,
    event_date timestamptz not null default now(),
    constraint chk_transactions_amount_positive check (amount != 0),
    constraint chk_transactions_currency_len_eq_3 check (char_length(currency) = 3)
);

create index idx_transactions_account_id on transactions (account_id);
commit;
