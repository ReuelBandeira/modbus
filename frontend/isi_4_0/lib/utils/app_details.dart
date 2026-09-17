enum AppDetails {
  name(text: 'ISI 4.0'),
  appVertion(text: ": 1.0.0 - 2023");

  final String text;
  const AppDetails({required this.text});
}

List<String> nameProtocols = ['Canopen', 'Ethernet/IP', 'Modbus'];

List<Map> flags = [
  {'id': 1, 'image': 'assets/flag_en.png', 'name': 'English', 'lang': 'en'},
  {'id': 2, 'image': 'assets/flag_es.png', 'name': 'Español', 'lang': 'es'},
  {'id': 3, 'image': 'assets/flag_pt.png', 'name': 'Português', 'lang': 'pt'}
];

List<Map<String, dynamic>> canopen = [
  {
    "nodeId": 279,
  }
];

List<Map<String, dynamic>> ethernetIP = [
  {
    "bitMemories": [
      // Flávio
      {"attribute": 0, "class": 849, "instance": 1, "name": "Y00"},
      {"attribute": 1, "class": 849, "instance": 1, "name": "Y01"},
      {"attribute": 2, "class": 849, "instance": 1, "name": "Y02"},
      {"attribute": 3, "class": 849, "instance": 1, "name": "Y03"}
    ],
    "wordMemories": [
      // Flávio
      {"attribute": 0, "class": 850, "instance": 2, "name": "D0"},
      {"attribute": 1, "class": 850, "instance": 2, "name": "D1"},
      {"attribute": 2, "class": 850, "instance": 2, "name": "D2"},
      {"attribute": 3, "class": 850, "instance": 2, "name": "D3"}
    ]
  }
];

List<Map<String, dynamic>> modbus = [
  {
    "bitMemories": [
      // Holden
      {"address": 40960, "name": "Y00"},
      {"address": 40961, "name": "Y01"},
      {"address": 40962, "name": "Y02"},
      {"address": 40963, "name": "Y03"},
    ],
    "slaveId": 1,
    "wordMemories": [
      // Holden
      {"format": 16, "address": 0, "name": "D0"},
      {"format": 16, "address": 1, "name": "D1"},
      {"format": 16, "address": 2, "name": "D2"},
      {"format": 16, "address": 3, "name": "D3"},
    ]
  }
];
