import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "globalData.dart";

class EditPasswordPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      title: "Account",
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  EditPassword createState() => EditPassword();
}

class EditPassword extends State<APICall> {
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  final TextEditingController oldPasswordController = TextEditingController();
  final TextEditingController newPasswordController = TextEditingController();
  final TextEditingController repeatNewPasswordController =
      TextEditingController();
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
      developer.log(
          'USER ERROR : $notSamePasswd : ${newPasswordController.text} ${repeatNewPasswordController.text}');
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
            'oldpassword': oldPasswordController.text,
            'password': newPasswordController.text,
          }),
        );

        if (response.statusCode == 200) {
          setState(() {
            result = "";
            developer.log("Password putData succes");

            Navigator.pop(Globaldata.myContext);
          });
        } else {
          final responseData = jsonDecode(response.body);
          developer.log(
              "Password putData FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
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
              (result.split(' ').contains("Invalid") == true &&
                  result.split(' ').contains("argument(s):") == true)) {
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
            EditPasswordTopBar(),
            SizedBox(width: 16),
            MyDivider(),
            Text(
              result,
              style: TextStyle(
                  fontSize: 20, fontWeight: FontWeight.bold, color: Colors.red),
            ),
            LeftWidget(
              child: AutoSizeText(
                "Current password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
                hint: 'Password',
                controller: oldPasswordController,
                obscureText: true,
                padding: EdgeInsets.symmetric(horizontal: 15)),
            Padding(padding: EdgeInsets.only(top: 30)),
            LeftWidget(
              child: AutoSizeText(
                "New password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
                hint: 'Password',
                controller: newPasswordController,
                obscureText: true,
                padding: EdgeInsets.symmetric(horizontal: 15)),
            Padding(padding: EdgeInsets.only(top: 20)),
            LeftWidget(
              child: AutoSizeText(
                "Confirm new password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
                hint: 'Password',
                controller: repeatNewPasswordController,
                obscureText: true,
                padding: EdgeInsets.symmetric(horizontal: 15)),
            LeftWidget(
              child: AutoSizeText(
                notSamePasswd,
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: const Color.fromARGB(255, 255, 0, 0),
                ),
                maxLines: 1,
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText(
                "Save",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: putData,
              size: Size(double.infinity, 55),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
    );
  }
}

class EditPasswordTopBar extends StatelessWidget {
  const EditPasswordTopBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(Icons.arrow_back, size: 32),
          ),
        ),
        Center(
          child: AutoSizeText(
            "Account",
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: Colors.black,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}
