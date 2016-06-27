#!/bin/bash

curl -XPUT -v 127.0.0.1:5000/v2/service_instances/123 -d '{"organization_guid": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
sleep 3
curl -XPUT -v 127.0.0.1:5000/v2/service_instances/123/service_bindings/456 -d '{"app_gui": "org-guid-here", "plan_id": "3C7BAF72-0DF4-420D-B365-B5CF09409C70","service_id":"04EB4D8F-15BF-43F2-B4DA-E7A243E21C83"}'
sleep 3
curl -XDELETE -v 127.0.0.1:5000/v2/service_instances/123
curl -XDELETE -v 127.0.0.1:5000/v2/service_instances/123/service_bindings/456
