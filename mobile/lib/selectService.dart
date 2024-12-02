import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "edit_account.dart";
import "globalData.dart";

class SelectServicePage extends StatelessWidget {

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      title: "Account",
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {

  @override
  SelectService createState() => SelectService();
}

class SelectService extends State<APICall> {
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');
  final TextEditingController oldPasswordController = TextEditingController();
  final TextEditingController newPasswordController = TextEditingController();
  final TextEditingController repeatNewPasswordController = TextEditingController();
  String result = '';
  String notSamePasswd = "";
  String apiRoute = 'user/modify-password';

  @override
  void dispose() {
    oldPasswordController.dispose();
    newPasswordController.dispose();
    super.dispose();
  }

  Future<void> putData() async {
    developer.log('SERVER_LINK : $apiUrl');
    if (repeatNewPasswordController.text != newPasswordController.text) {
      setState(() {
        notSamePasswd = "Password doesn't match";
      });
      developer.log('USER ERROR : $notSamePasswd : ${newPasswordController.text} ${repeatNewPasswordController.text}');
    } else {
      setState(() {
        notSamePasswd = "";
      });
      try {
        developer.log('SERVER_LINK + redirection : ${apiUrl + apiRoute}');
        final response = await http.put(
          Uri.parse(apiUrl + apiRoute),
          headers: <String, String>{
            'Content-Type': 'application/json; charset=UTF-8',
            'Cookie': Globaldata.JWToken,
          },
          body: jsonEncode(<String, dynamic>{
            'oldpassword' : oldPasswordController.text,
            'password': newPasswordController.text,
          }),
        );

        if (response.statusCode == 200) {
          setState(() {
            result = "";
            developer.log("Password putData succes");

            Navigator.push(
              Globaldata.myContext,
              MaterialPageRoute(builder: (context) => EditAccountPage()),
            );
          });
        } else {
          final responseData = jsonDecode(response.body);
          developer.log("Password putData FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
          setState(() {
            result = '${responseData['error']}';
          });
        }
      } catch (e) {
        developer.log("Password putData ERROR: $e", level: 1);
        setState(() {
          result = '$e';
          developer.log(result);
          if (result.split(' ').contains("SocketException:") == true ||
             (result.split(' ').contains("Invalid") == true && result.split(' ').contains("argument(s):") == true)) {
            result = Globaldata.ErrorServerConnection;
          }
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectServiceTopBar(title: "title"),

            SizedBox(width: 16),
            MyDivider(),
            Text(result, style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.red),),
          ],
        ),
      ),
    );
  }
}

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
  });

  final String title;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {Navigator.pop(Globaldata.myContext);},
            icon: Icon(Icons.arrow_back, size: 32),
          ),
        ),
        Center(
          child: AutoSizeText(title, style: TextStyle( fontSize: 30, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),
        ),
      ],
    );
  }
}
