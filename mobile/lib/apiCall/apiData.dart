import 'package:flutter_dotenv/flutter_dotenv.dart';

class ApiData {
  static String apiUrl = dotenv.get('SERVER_LINK',
        fallback: 'Cannot find .env variable "SERVER_LINK"');
}