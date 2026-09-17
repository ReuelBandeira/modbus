import 'dart:convert';

class JsonData {
  String address;
  String port;
  String name;
  String? protocolConnection;
  String readingTime;
  List<String> topics;
  List<Map<String, dynamic>> data;

  JsonData({
    required this.address,
    required this.port,
    required this.name,
    required this.protocolConnection,
    required this.readingTime,
    required this.topics,
    required this.data,
  });

  Map<String, dynamic> toJson() {
    String topicEncode = jsonEncode(topics);
    String dataEncode = jsonEncode(data);

    return {
      'address': address,
      'port': port,
      'name': name,
      'protocol': protocolConnection,
      'readingTime': readingTime,
      'topics': topicEncode,
      'data': dataEncode,
    };
  }
}
