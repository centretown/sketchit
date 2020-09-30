# curl -H "login:testing" -H "password:test" -X POST -d '{"greeting":"foo"}' 'http://dragon:7778/1/ping'
# curl -H "login:testing" -H "password:test" -X GET 'http://dragon:7778/v1/sectors/work/devices/esp32-01'
# curl -H "login:testing" -H "password:test" -X GET  'http://dragon:7778/v1/sectors//devices'
echo ""
echo "List"
curl -H "login:testing" -H "password:test" -X GET  'http://dragon:7778/v1/sectors/cottage/devices'
echo ""
echo "Get"
curl -H "login:testing" -H "password:test" -X GET  'http://dragon:7778/v1/sectors/cottage/devices/esp32-01'
echo ""
echo "Create 02"
curl -H "login:testing" -H "password:test" -X POST -d \
    '{"sector":"cottage","label":"esp32-02","toolkit":"ESP32","pins":[{"id":2,"label":"LED","purpose":"activity indicator"},{"id":5,"label":"TX","purpose":"soft serial transmitter"},{"id":6,"label":"RX","purpose":"soft serial receiver"}],"ip":"192.168.1.202:8080","port":"esp32-02"}' \
    'http://dragon:7778/v1/sectors/cottage/devices'
echo ""
echo "Get 02"
curl -H "login:testing" -H "password:test" -X GET 'http://dragon:7778/v1/sectors/cottage/devices/esp32-02'
echo ""
echo "Update 02"
curl -H "login:testing" -H "password:test" -X PATCH -d \
    '{"sector":"cottage","label":"esp32-02","toolkit":"WHOLE WEED BREAD","pins":[{"id":2,"label":"LED","purpose":"activity indicator"},{"id":5,"label":"TX","purpose":"soft serial transmitter"},{"id":6,"label":"RX","purpose":"soft serial receiver"}],"ip":"192.168.1.202:8080","port":"esp32-02"}' \
    'http://dragon:7778/v1/sectors/cottage/devices/esp32-02'
echo ""
echo "Get 02"
curl -H "login:testing" -H "password:test" -X GET 'http://dragon:7778/v1/sectors/cottage/devices/esp32-02'
echo ""
echo "Delete 02"
curl -H "login:testing" -H "password:test" -X DELETE 'http://dragon:7778/v1/sectors/cottage/devices/esp32-02'
echo ""
echo ""

echo "Create test data"
echo "Create 02"
curl -H "login:testing" -H "password:test" -X POST -d \
    '{"sector":"cottage","label":"esp32-02","toolkit":"ESP32","pins":[{"id":2,"label":"LED","purpose":"activity indicator"},{"id":5,"label":"TX","purpose":"soft serial transmitter"},{"id":6,"label":"RX","purpose":"soft serial receiver"}],"ip":"192.168.1.202:8080","port":"esp32-02"}' \
    'http://dragon:7778/v1/sectors/cottage/devices'

echo "Create 03"
curl -H "login:testing" -H "password:test" -X POST -d \
    '{"sector":"cottage","label":"esp32-03","toolkit":"ESP32","pins":[{"id":2,"label":"LED","purpose":"activity indicator"},{"id":5,"label":"TX","purpose":"soft serial transmitter"},{"id":6,"label":"RX","purpose":"soft serial receiver"}],"ip":"192.168.1.202:8080","port":"esp32-02"}' \
    'http://dragon:7778/v1/sectors/cottage/devices'

echo "Create 04"
curl -H "login:testing" -H "password:test" -X POST -d \
    '{"sector":"cottage","label":"esp32-04","toolkit":"ESP32","pins":[{"id":2,"label":"LED","purpose":"activity indicator"},{"id":5,"label":"TX","purpose":"soft serial transmitter"},{"id":6,"label":"RX","purpose":"soft serial receiver"}],"ip":"192.168.1.202:8080","port":"esp32-02"}' \
    'http://dragon:7778/v1/sectors/cottage/devices'
