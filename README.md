# OpsAgent

Local-first Linux AI ops agent. Tek kurulum komutuyla sunucuya kurulur, kendi SQLite veritabanını kullanır, web paneline sahiptir.

**Güvenlik ilkesi:** Read-only varsayılan. Sistem değişikliği gerektiren her işlem açık onay gerektirir.

## Tek Kurulum

```bash
curl -fsSL https://opsagent.dev/install.sh | sudo bash
```

Kurulum sonrası systemd servisi aktif olur. Dashboard: `http://127.0.0.1:8787`

## Uzak Erişim

Dashboard localhost'ta çalışır. Uzaktan erişim için SSH tunnel:

```bash
ssh -L 8787:localhost:8787 root@SUNUCU_IP
```

Tarayıcıda: `http://localhost:8787`

## Kurulum Adımları

```
1. Root yetkisi kontrolü
2. OS/architecture tespiti (Linux x86_64 veya arm64)
3. Binary indirme veya yerel build kullanımı
4. /usr/local/bin/opsagent kurulumu
5. /etc/opsagent dizini oluşturma
6. config.yaml oluşturma
7. /var/lib/opsagent dizini (SQLite veritabanı)
8. /var/log/opsagent dizini (loglar)
9. systemd service dosyası oluşturma
10. Servis başlatma ve enable etme
```

Kurulum başarılı olursa:

```
OpsAgent installed successfully.

Service:
active

Dashboard:
http://127.0.0.1:8787

Your server data stays on this machine.
No write operation will run without approval.
```

## Manuel Build

```bash
# UI build
cd web/dashboard && npm install && npm run build && cd ../..

# Go binary build
go build -o opsagent ./cmd/opsagent
```

## Manuel Kurulum (SSH ile erişim yoksa)

```bash
# Binary'yi sunucuya kopyala
scp opsagent root@SUNUCU_IP:/tmp/

# Sunucuda çalıştır
ssh root@SUNUCU_IP
chmod +x /tmp/opsagent
sudo /tmp/opsagent serve --config /etc/opsagent/config.yaml
```

## Kaldırma

```bash
curl -fsSL https://opsagent.dev/uninstall.sh | sudo bash
```

## Yapı

```
/usr/local/bin/opsagent          # Ana binary
/etc/opsagent/config.yaml       # Config dosyası
/var/lib/opsagent/opsagent.db    # SQLite veritabanı
/var/log/opsagent/opsagent.log  # Log dosyası
/etc/systemd/system/opsagent.service
```

## İlk Kullanım

1. Dashboard'u aç: `http://127.0.0.1:8787`
2. İlk açılışta setup wizard çalışır
3. Admin hesabı oluştur
4. AI provider istersen API key gir (opsiyonel)
5. Agent otomatik metrik toplamaya başlar

## Özellikler

- **Monitoring:** CPU, RAM, disk, load average, processler, portlar, servisler
- **Alerting:** Disk/RAM/CPU threshold bazlı uyarılar
- **Assistant:** Doğal dilde sistem sorgulama
- **Plan + Approval:** Sistem değişikliği gerektiren işlemlerde onay akışı
- **Audit Log:** Tüm işlemlerin kaydı

## AI Olmadan Kullanım

AI provider yapılandırılmamışsa agent template bazlı planlar üretir. Disk sorguları `df -h`/`du` komutlarıyla, diğer sorgular `uptime`/`free -m` ile yanıtlanır.

## Güvenlik

- Dashboard varsayılan olarak sadece 127.0.0.1'de çalışır
- Write komutlar onaysız çalışmaz
- Blocked komutlar (rm -rf /, mkfs, vb.) asla çalıştırılamaz
- Tüm işlemler audit log'da kayıtlı
