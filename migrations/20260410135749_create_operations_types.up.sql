begin;

create table if not exists operations_types (
    operation_type_id bigint primary key,
    description text not null,
    sign_multiplier smallint not null,
    created_at timestamptz not null default now(),
    constraint chk_operations_types_sign_multiplier check (sign_multiplier in (-1, 1)),
    constraint chk_operations_types_description_not_empty check (btrim(description) != ''),
    constraint chk_opeartions_types_max_len check (char_length(description) <= 30)
);

insert into operations_types (operation_type_id, description, sign_multiplier) values
(1, 'PURCHASE', -1),
(2, 'INSTALLMENT PURCHASE', -1),
(3, 'WITHDRAWAL', -1),
(4, 'PAYMENT', 1);

commit;
