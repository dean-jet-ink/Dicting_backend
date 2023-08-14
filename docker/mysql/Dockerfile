FROM mysql:8.0

# 設定ファイルコピー
COPY ./my.cnf /etc/mysql/conf.d/my.cnf

# コンテナ内に事前定義されているユーザー"mysql"に権限を限定
RUN chown -R mysql /var/lib/mysql && \
chgrp -R mysql /var/lib/mysql