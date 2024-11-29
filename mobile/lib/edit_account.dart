import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:auto_size_text/auto_size_text.dart';

import "custom_widget.dart";
import "explore.dart";
import "footer_page.dart";
import "editPassword.dart";

class EditAccount extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      title: "Account",
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  _APICallState createState() => _APICallState();
}

class _APICallState extends State<APICall> {
  final String apiUrl = dotenv.get('SERVER_LINK', fallback: 'Cannot find .env variable "SERVER_LINK"');
  final TextEditingController passwordController = TextEditingController();
  String result = '';
  String apiRoute = 'modify-password';

  @override
  void dispose() {
    passwordController.dispose();
    super.dispose();
  }

  Future<void> _postData() async {
    developer.log('SERVER_LINK : $apiUrl');
    try {
      developer.log('SERVER_LINK + redirection : ${apiUrl + apiRoute}');
      final response = await http.post(
        Uri.parse(apiUrl + apiRoute),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
        },
        body: jsonEncode(<String, dynamic>{
          'password': passwordController.text,
        }),
      );

      if (response.statusCode == 200) {
        setState(() {
          result = "";

          Navigator.pop(context);
          Navigator.push(
            context,
            MaterialPageRoute(builder: (context) => ExplorePage()),
          );
        });
      } else {
        final responseData = jsonDecode(response.body);
        setState(() {
          result = '${responseData['error']}';
        });
      }
    } catch (e) {
      setState(() {
        result = '$e';
        if (result.split(' ').contains("SocketException:") == true ||
           (result.split(' ').contains("Invalid") == true && result.split(' ').contains("argument(s):") == true)) {
          debugPrint(result);
          result = "Could not find the server\nPlease try again later";
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            AutoSizeText("Account", style: TextStyle (fontSize: 30, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),
            SizedBox(width: 16),
            MyDivider(),
            Padding(padding: EdgeInsets.only(top: 30)),

            LeftWidget(child: AutoSizeText("Email", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),),

            CustomClassicForm(hint:'Email', padding: EdgeInsets.symmetric(horizontal: 15),),
            Padding (padding: EdgeInsets.only(top: 20)),
            LeftWidget(child: AutoSizeText("Password", style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold, color: Colors.black,), maxLines: 1,),),
            CustomClassicForm(hint: 'Password', controller: passwordController, obscureText: true, padding: EdgeInsets.symmetric(horizontal: 15)),
            LeftWidget(
              child: IntrinsicWidth(
                child: TextButton(
                  onPressed: () {Navigator.push(context, MaterialPageRoute(builder: (context) => EditPassword()),);},
                  style: TextButton.styleFrom(
                    padding: EdgeInsets.zero,
                  ),
                  child: const Text(
                    "Change password",
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: Color.fromARGB(255, 0, 160, 255),
                    ),
                  ),
                ),
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 20)),
            CustomButton(
              backgroundColor:Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText("Save", style: TextStyle(fontSize: 30, fontWeight: FontWeight.bold, color: Colors.white),),
              onPressed: () {},
              size: Size(double.infinity, 55),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 4),
    );
  }
}

class FormWithTitle extends StatelessWidget {
  const FormWithTitle({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Column (
      children: [
        AutoSizeText("All", style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.black), maxLines: 1, textAlign: TextAlign.right,),
        CustomClassicForm(hint:'Email'),
      ],
    );
  }
}
